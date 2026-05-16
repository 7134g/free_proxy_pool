## 1. 依赖与配置

- [x] 1.1 添加 `github.com/redis/go-redis/v9` 依赖到 go.mod
- [x] 1.2 在 config 结构体中新增 `redis.enabled` 字段
- [x] 1.3 更新 config.yaml 示例配置添加 `enabled` 项

## 2. Redis 连接与存储层

- [x] 2.1 实现 `config.InitRedis()`：建立 Redis 客户端连接，PING 检测可用性
- [x] 2.2 新建 `crawler/redis_store.go`：封装 Redis ZSET 操作（ZADD、ZINCRBY、ZREM、ZREVRANGEBYSCORE）
- [x] 2.3 Redis 操作失败时 log 警告而不阻断主流程

## 3. Store 集成 Redis 持久化

- [x] 3.1 修改 `Store.add()`：写入内存后同步 ZADD 到 Redis
- [x] 3.2 修改 `Store.inc()`：Score 变更后同步 ZINCRBY 到 Redis
- [x] 3.3 修改 `Store.dnc()`：Score 变更后同步 ZINCRBY；删除时同步 ZREM
- [x] 3.4 修改 `Store.Del()`：删除时同步 ZREM
- [x] 3.5 新增启动恢复函数：从 Redis ZSET 加载所有 score > 0 的代理到内存 Store

## 4. Docker 容器化

- [x] 4.1 编写多阶段构建 Dockerfile（golang:1.23-alpine 编译 → alpine:3.20 运行）
- [x] 4.2 编写 .dockerignore 排除不必要文件
- [x] 4.3 验证 `docker build` 可正常构建镜像（本地无 Docker，Dockerfile 已编写完成，需在有 Docker 环境验证）
