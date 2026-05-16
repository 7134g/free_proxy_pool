# CLAUDE.md

本文件为 Claude Code（claude.ai/code）在操作此代码库时提供指引。

## 构建与运行

```bash
go mod tidy
go build .
./free_proxy_pool --config ./config.yaml
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行单个爬虫测试
go test -v -run TestCrawler_crawlDaiLi66 ./crawler/cell/

# 运行代理测试
go test -v -run TestRunTestProxy ./crawler/
```

## 项目架构

免费代理池项目，从多个公开网站抓取免费代理 IP，经过可用性验证后按评分维护在内存中，通过 HTTP API 和转发代理两种方式对外提供服务。

### 目录结构

- `main.go` — 入口，启动三个并发模块：API 服务、Martian 代理、爬虫
- `config/` — YAML 配置解析（Redis 配置预留，尚未实现）
- `crawler/` — 核心爬虫逻辑
  - `cell/spider.go` — spider 接口定义：`name()`、`genSeek()`、`run()`、`parse()`
  - `cell/crawl_*.go` — 各网站爬虫实现（66ip、ip3366、快代理、proxy11）
  - `cell/register.go` — 爬虫注册与并发执行
  - `run.go` — cron 定时抓取/测试，monitor 事件循环处理代理验证结果
  - `store.go` — 内存代理存储（`Store`），自旋锁保护的 concurrent map，支持基于`Score`的排序和增减
  - `tester.go` — 代理可用性测试（通过 test_urls 检测）
  - `init.go` — 全局状态：`CacheProxyData`、`TaskPool`（协程池）、`ProxyFinishChannel`
- `serve/` — HTTP 服务
  - `run.go` — Gin 引擎启动
  - `route.go` — 路由注册（`/max`、`/random`、`/list`、`/count`、`/config`、`/useless`）
  - `handler.go` — 请求处理器
  - `martian.go` — 基于 google/martian 的转发代理，定时从代理池选取最优代理
- `util/` — 工具包
  - `pool/` — 自实现 goroutine 协程池（Worker 模式，支持容量限制、空闲回收）
  - `proxy/` — TCP 级正向代理实现（HTTP CONNECT + 普通 HTTP 转发）
  - `cas/` — 自旋锁实现
  - `xhttp/` — 带代理设置的 HTTP 客户端封装

### 数据流

1. **爬虫**：cron 触发各 spider 并行抓取 → 解析 HTML 提取 `ip:port` → 发送到 `ProxyChannel`
2. **验证**：monitor 从 `ProxyChannel` 取出代理 → 提交到 `TaskPool` → 通过 `test_urls` 检测可用性 → 结果发送到 `ProxyFinishChannel`
3. **评分**：验证通过的代理 `Score++`，失败的 `Score--`，`Score <= 0` 时从池中移除。按分数降序排列
4. **对外服务**：
   - API 模式：`/max` 返回 Top10 中随机一个，`/random` 全池随机，`/list` 返回 Top10 列表
   - 代理模式：martian 定时从池中取最优代理，转发所有用户请求

### 配置（config.yaml）

关键配置项：`pool_cap`（代理池上限）、`flash_score`（新鲜度）、`test_time`/`crawler_time`（cron 表达式）、`martian.mode`（`random` 或 `max`）、`test_urls`（验证代理的测试链接）

### 设计特点

- 代理评分机制：成功加分、失败减分，低于阈值自动淘汰，保持池中代理质量
- 自实现协程池：控制并发爬取和测试的 goroutine 数量
- 基于自旋锁的并发安全 Store，非标准 sync.Map
- 支持 `max`（取评分最高）和 `random`（随机）两种代理选择策略
