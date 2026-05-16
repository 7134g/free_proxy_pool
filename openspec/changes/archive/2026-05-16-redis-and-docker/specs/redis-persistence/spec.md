## ADDED Requirements

### Requirement: Redis persistence for proxy data
The system SHALL persist proxy data to Redis when `redis.enabled` is true, and SHALL recover proxy data from Redis on startup.

#### Scenario: Proxy added to pool is persisted to Redis
- **WHEN** a new proxy passes validation and is added to the pool (Score set to `flash_score`)
- **THEN** the proxy SHALL be written to the Redis ZSET with its `ip:port` as member and Score as score

#### Scenario: Proxy score increased is synced to Redis
- **WHEN** a proxy passes a subsequent validation test and its Score is incremented
- **THEN** the Redis ZSET score for that proxy SHALL be incremented accordingly

#### Scenario: Proxy score decreased is synced to Redis
- **WHEN** a proxy fails a validation test and its Score is decremented
- **THEN** the Redis ZSET score for that proxy SHALL be decremented accordingly

#### Scenario: Proxy removed from pool is deleted from Redis
- **WHEN** a proxy's Score drops to 0 or below and is removed from the pool
- **THEN** the proxy SHALL be removed from the Redis ZSET

#### Scenario: Pool recovery from Redis on startup
- **WHEN** the service starts with `redis.enabled` set to true
- **THEN** all proxies with Score > 0 SHALL be loaded from Redis ZSET into the in-memory Store

#### Scenario: Graceful degradation when Redis is unavailable
- **WHEN** Redis is unreachable or `redis.enabled` is false
- **THEN** the system SHALL continue operating normally using in-memory Store only, logging a warning for failed Redis operations

### Requirement: Redis connection management
The system SHALL establish and manage a Redis client connection based on configuration.

#### Scenario: Redis connection initialized from config
- **WHEN** the service starts with `redis.enabled` true, `redis.url` set, and optional `redis.password`
- **THEN** a Redis client SHALL be created and connectivity verified via PING

#### Scenario: Startup fails fast on bad Redis config
- **WHEN** the service starts with `redis.enabled` true but the Redis connection fails
- **THEN** the service SHALL log the error and continue running in memory-only mode

### Requirement: Docker containerization
The system SHALL provide a Dockerfile for building a container image.

#### Scenario: Multi-stage Docker build
- **WHEN** `docker build` is executed in the project root
- **THEN** a minimal Alpine-based image SHALL be produced containing the compiled binary and default config

#### Scenario: Container runs the proxy pool service
- **WHEN** the Docker container is started with appropriate port mappings
- **THEN** the proxy pool service SHALL start and serve HTTP API and forward proxy as configured
