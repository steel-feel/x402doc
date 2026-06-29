# TakeDoc — x402 Paid Documents Service

A Go service that gates document access behind x402 blockchain payments on Tempo testnet, using hexagonal architecture with Echo, SQLite, and gRPC.

## Design Decisions (User-Confirmed)

| Decision | Choice |
|---|---|
| gRPC role | TakeDoc exposes a gRPC **server** — other services call in to get document details |
| Document storage | SQLite (content stored as text in DB) |
| Target network | Tempo **Testnet** |
| x402 facilitator | **Self-hosted** (embedded facilitator package within TakeDoc) |
| Dependency injection | **Manual constructor injection** (no DI framework) |

---

## Target Project Structure

```
prac/
├── cmd/server/main.go
├── internal/
│   ├── domain/          # Entities: Document, AccessLog, HealthStatus
│   ├── port/            # Interfaces: incoming (services) + outgoing (repos)
│   ├── service/         # Business logic implementations
│   └── adapter/
│       ├── primary/
│       │   ├── http/    # Echo handlers + middleware (x402, access log, observability)
│       │   └── grpc/    # gRPC server for document queries
│       └── secondary/
│           ├── sqlite/  # DB repos + migrations
│           └── facilitator/  # Self-hosted x402 facilitator
├── proto/document/v1/   # Protobuf definitions
├── api/generated/       # Generated Go code
├── migrations/          # SQL migration files
├── spec/main.tsp        # TypeSpec API definition
└── smoke_test.go        # Integration tests
```

---

## Phases

### Phase 1: Foundation & Hexagonal Architecture
- Restructure codebase into hexagonal layers (domain → ports → services → adapters)
- SQLite setup with migrations and health check
- Complete health endpoint (status, db_status, uptime)

### Phase 2: Document Domain & REST API
- Document entity, repository, service
- `GET /api/v1/documents/:id` endpoint (no payment gating yet)
- Seed data in migrations

### Phase 3: x402 Payment Integration
- Echo middleware for x402 flow (402 response → payment → verification)
- Self-hosted facilitator for Tempo testnet
- Payment gating on document endpoints

### Phase 4: gRPC Server
- Proto definitions for DocumentService (GetDocument, ListDocuments)
- gRPC server running alongside HTTP on separate port
- Shared domain service between HTTP and gRPC

### Phase 5: Access Logging & Observability
- Log document access (blockchain address, IP, user agent, tx hash) to SQLite
- Structured logging with `log/slog`
- Prometheus metrics + OpenTelemetry tracing

### Phase 6: Smoke Testing & Polish
- Integration tests for all endpoints (health, documents, x402 flow, gRPC)
- Makefile targets: test, lint, build, proto-gen
- TypeSpec spec updates

---

See the full detailed plan in the implementation_plan.md artifact for file-level details and verification steps.
