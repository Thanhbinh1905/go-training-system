# SETA Golang Training - Microservices System

## Overview

A production-style microservices system in Go for managing Users, Teams, and Assets (folders/notes). It demonstrates:
- Clean service boundaries (User, Team, Asset)
- PostgreSQL + GORM for persistence
- Redis for caching and real-time ACL lookup
- gRPC between services + GraphQL (User) + REST (Team/Asset)
- Kafka for domain events (team/activity, asset/changes)
- Observability with Prometheus + Grafana + Loki + Promtail
- Graceful shutdown

## Services

- User Service (GraphQL + REST)
  - Auth, JWT issuance/verification, user profile and roles (MANAGER/MEMBER)
  - Exposes gRPC API for token verification and user lookups
  - Port: 8081 (HTTP), 50051 (gRPC)

- Team Service (REST + gRPC)
  - Create team, add/remove managers and members
  - Consumes User gRPC to validate users
  - Emits Kafka events for team changes
  - Port: 8082 (HTTP), 50052 (gRPC)

- Asset Service (REST)
  - Create/Update/Delete folders and notes
  - Share/Unshare folders/notes with access levels (read/write)
  - Real-time ACL cache in Redis: asset:{assetId}:acl → { userId: accessType }
  - Validates permissions from Redis before DB fallbacks
  - Emits Kafka events for asset changes
  - Port: 8083 (HTTP)

## Architecture

```text
Client ──► (GraphQL/REST)
   ├─► user-service (GraphQL + gRPC)
   ├─► team-service (REST + gRPC)
   └─► asset-service (REST)

user-service ◄── gRPC ──► team-service
asset-service ◄── gRPC ──► user-service, team-service

Kafka: team.activity, asset.changes
Redis: caching, real-time ACL
PostgreSQL: per-service DB
Prometheus/Grafana + Loki/Promtail: metrics/logs
```

## Tech Stack

- Go 1.24, Gin, gqlgen, gRPC, GORM
- PostgreSQL (per-service DB): user, team, asset
- Redis (caching + ACL): redis:7.2
- Kafka (segmentio/kafka-go) + Zookeeper
- Prometheus, Grafana, Loki, Promtail
- Go workspaces (go.work)

## Kafka Integration

- Topics
  - `team.activity`: TEAM_CREATED, MEMBER_ADDED, MEMBER_REMOVED, MANAGER_ADDED, MANAGER_REMOVED
  - `asset.changes`: FOLDER_CREATED/UPDATED/DELETED, NOTE_CREATED/UPDATED/DELETED, FOLDER_SHARED/UNSHARED, NOTE_SHARED/UNSHARED

- Producers
  - Team Service publishes team activity events
  - Asset Service publishes asset change and share/unshare events

- Shared package `shared/kafka`
  - `events.go`: event types and payloads
  - `producer.go`: publish helpers (team, asset, asset share)
  - `consumer.go` + `consumer_service.go`: consuming patterns and sample handlers

See `KAFKA_INTEGRATION.md` for details and Taskfile commands.

## Real-Time ACL (Redis)

- Key format: `asset:{assetId}:acl` (Redis Hash)
  - Field: `userId`
  - Value: `read` or `write`

- Updated on share/unshare
  - On ShareFolder/ShareNote: set multiple user access in the hash
  - On RevokeFolderShare/RevokeNoteShare: remove user’s field from the hash

- Permission validation flow
  - First check Redis ACL
  - If not found, fallback to DB checks (NoteShare/FolderShare and owner logic)

Relevant code:
- `services/asset-service/internal/repository/asset_cache.go` (ACL helpers)
- `services/asset-service/internal/repository/asset_repository.go` (write-through cache on share/unshare)
- `services/asset-service/internal/service/asset_service.go` (permission checks prioritize Redis)

## Graceful Shutdown

- Shared utilities in `shared/graceful`
- Each service starts HTTP/gRPC servers in goroutines and listens for SIGINT/SIGTERM
- Shutdown completes in a bounded timeout (configurable per service)

See `GRACEFUL_SHUTDOWN.md` for usage.

## Project Structure

```text
.
├── services/
│   ├── user-service/
│   │   ├── cmd/                     # entrypoint(s)
│   │   ├── internal/                # app, handler (graph/http/grpc), service, repo, model
│   │   ├── pb/                      # generated from shared/proto
│   │   ├── config/
│   │   └── Dockerfile
│   ├── team-service/
│   │   ├── internal/                # gRPC + REST handlers, service, repo, middleware
│   │   ├── pb/
│   │   ├── config/
│   │   └── Dockerfile
│   └── asset-service/
│       ├── internal/                # REST handlers, service, repo (DB+cache), model
│       ├── pb/
│       ├── config/
│       └── Dockerfile
├── shared/
│   ├── proto/                       # user, team protobufs
│   ├── kafka/                       # events, producer, consumer
│   ├── db/ (postgres, redis)
│   ├── logger/
│   ├── graceful/
│   └── errors, apperror, contextkey, httpresponse
├── migration/                       # GORM migrations for all DBs
├── docker-compose.yml               # services + infra
├── Taskfile.yaml                    # dev helpers
├── KAFKA_INTEGRATION.md
└── GRACEFUL_SHUTDOWN.md
```

## Run with Docker

Requirements: Docker, docker-compose. Set environment variables (see docker-compose references):
- POSTGRES_USER, POSTGRES_PASSWORD
- JWT_SECRET
- PRODUCTION (true/false)

Start infra + services:
```bash
docker-compose up -d --build
```
Services:
- user-service: http://localhost:8081, gRPC 50051
- team-service: http://localhost:8082, gRPC 50052
- asset-service: http://localhost:8083
- Kafka: localhost:9092
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100

## Development

Using Taskfile:
```bash
# Protobuf codegen
task proto:user
task proto:team
task proto:asset

# gqlgen (user-service GraphQL)
task user:gqlgen

# DB shells (inside Docker)
task pqsl:user
task pqsl:team
task pqsl:asset

# Redis shell
task redis-cli

# Kafka helper commands
task kafka-topics
task kafka-console-consumer -- team.activity
task kafka-console-consumer -- asset.changes
task kafka-console-producer -- team.activity
```

Run locally with Air (hot reload) using service Dockerfiles/volumes or your preferred method.

## API Highlights

- User Service (GraphQL + REST)
  - GraphQL endpoint: `/graphql` (+ Playground on `/`)
  - REST: `/api/v1/health`, `/api/v1/users` (import from file)

- Team Service (REST + gRPC)
  - REST (manager-protected): `/api/v1/teams/...`
  - gRPC: team definitions in `shared/proto/team/team.proto`

- Asset Service (REST)
  - REST: `/api/v1/assets/...`
  - Real-time ACL checked from Redis before DB

## Observability

- Logs: Promtail ships logs from service volumes to Loki
- Metrics: Prometheus scrapes; Grafana dashboards available (login from env)

## Notes

- This repo uses Go workspaces (`go.work`) to develop shared modules alongside services.
- Kafka producers/consumers use segmentio/kafka-go.
- ACL cache is the first line of authorization checks; DB is the source of truth.

## License

MIT – Free to learn, use, and modify.

---

## 🙌 Credits

Project built as part of **SETA Training Program** to learn production-grade Golang microservice architecture.

---

