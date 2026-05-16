## Context

当前代理池的代理数据完全存储在内存中（`crawler.Store`），服务重启即丢失。config.yaml 已预留 `redis` 配置段，`config.InitRedis()` 仅有 TODO 占位。项目无 Docker 化支持，本地/服务器部署需手动管理 Go 环境和编译。

## Goals / Non-Goals

**Goals:**
- 代理数据持久化到 Redis，服务重启后自动恢复，避免冷启动重新抓取
- 内存 Store 保留为热数据层，Redis 为持久化层；评分变更实时同步到 Redis
- Redis 连接失败时优雅降级，不影响核心抓取/验证/服务功能
- 提供 Dockerfile（多阶段构建）和 .dockerignore，支持容器化部署
- Redis 可通过配置开关完全关闭，回退到现有纯内存模式

**Non-Goals:**
- 不实现 Redis Sentinel/Cluster 高可用方案（单机 Redis 足够）
- 不替换现有内存 Store 架构，Redis 仅作为持久化补充
- 不提供 docker-compose（用户可以自行组合）

## Decisions

### 1. Redis 数据模型：使用 Sorted Set（ZSET）

选择 ZSET 而非 Hash + 手动排序：
- ZSET 天然按 Score 排序，与代理评分机制完全契合
- `ZADD` 写入，`ZREVRANGE` 读取 Top N，一个命令完成
- 单个 key 存储所有代理，member 为 `ip:port` 字符串，score 为评分值
- 备选方案 Hash：需额外维护排序列表，增加复杂度

### 2. 同步策略：写穿透（Write-Through）

- **写入**：代理评分变更时同步写 Redis（`ZADD`/`ZINCRBY`），不增加后台同步循环
- **恢复**：启动时 `ZREVRANGEBYSCORE` 读取所有 score > 0 的代理恢复到内存 Store
- **删除**：score <= 0 时 `ZREM` 移除
- 备选方案 Write-Back：延迟同步可减少 Redis 写入，但重启可能丢失最近变更，且增加代码复杂度

### 3. Redis 不可用时的降级策略

- Redis 操作失败仅 log 警告，不返回 error 阻断主流程
- 封装 `RedisStore` 接口，内部捕获连接错误
- 内存 Store 保持独立运作，不受 Redis 状态影响

### 4. Dockerfile：多阶段构建 + Alpine

- 第一阶段：`golang:1.23-alpine` 编译
- 第二阶段：`alpine:3.20` 运行，仅包含二进制和 config.yaml
- 使用 Alpine 减小镜像体积（预计 < 20MB）

### 5. 配置开关：`redis.enabled`

- `config.yaml` 新增 `enabled` 字段，默认 `false`（向后兼容）
- `enabled: true` 时初始化 Redis 连接并启用持久化
- 不传或 `false` 时行为与现有版本完全一致

## Risks / Trade-offs

- **Redis 写入延迟影响代理验证吞吐** → Redis 操作异步化，或用 `Pipeline` 批量写入；实际上单次 ZADD 延迟 < 1ms，影响可忽略
- **Redis 数据与内存不一致** → 内存 Store 始终是权威数据源，Redis 仅在 score 变更时同步，不会出现 Redis 覆盖内存的情况
- **Docker 容器内无法访问宿主机 Redis（127.0.0.1）** → config.yaml 中 redis.url 需配置为实际 Redis 地址；Docker 部署时通过环境变量或挂载配置覆盖

## Migration Plan

1. 添加 Redis 依赖，实现 `RedisStore`
2. 在 `config.InitRedis()` 中初始化连接
3. 修改 `crawler.Store` 在 inc/dnc/add/del 时同步写 Redis
4. 启动时从 Redis 恢复数据到 Store
5. 编写 Dockerfile 和 .dockerignore
6. 测试：纯内存模式（redis.enabled=false）行为不变；Redis 模式数据持久化和恢复

## Open Questions

- Redis 密码是否支持（已预留配置字段，本次实现）
