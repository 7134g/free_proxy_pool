## Why

当前代理池仅使用内存存储，服务重启后所有已抓取的代理数据全部丢失，冷启动需要重新抓取和验证，恢复时间长。同时项目缺少容器化支持，部署和分发不便。引入 Redis 持久化解决数据丢失问题，添加 Docker 支持简化部署流程。

## What Changes

- 新增 Redis 存储层，启动时从 Redis 恢复代理数据，代理评分变更时同步写入 Redis
- 内存 Store 保留作为热数据层，Redis 作为持久化层，二者协同工作
- 实现 `config.InitRedis()`，建立 Redis 连接池并检测可用性
- 新增 `Dockerfile`（多阶段构建）和 `.dockerignore`，支持容器化部署
- 新增 `redis.enabled` 配置项，允许关闭 Redis 回退到纯内存模式

## Capabilities

### New Capabilities

- `redis-persistence`: 代理数据持久化到 Redis，服务重启后可恢复；评分变更实时同步；支持通过配置开关启用/关闭

### Modified Capabilities

<!-- 无现有 spec 需要修改 -->

## Impact

- 依赖新增：`github.com/redis/go-redis/v9`
- 配置新增：`redis.enabled` 字段
- 影响文件：`config/yaml.go`、`config.yaml`、`crawler/store.go`、`crawler/init.go`、`go.mod`
- 新增文件：`redis_store.go`（Redis 存储操作）、`Dockerfile`、`.dockerignore`
