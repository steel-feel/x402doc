# Hexagonal Document Service with REST, gRPC & JWT Auth

A Go microservice designed around Hexagonal Architecture (Ports & Adapters) providing REST and gRPC interfaces for document management, payment gating, and market data retrieval.

---

## Architecture & Tech Stack

- **Go**: Version `1.26+`.
- **TypeSpec**: Used to define API contracts, automatically generating OpenAPI schemas and Protobuf specifications.
- **REST API**: Echo web framework. Gated document access via `X-402` payment validation and protected document creation using JWT authentication.
- **gRPC API**: Protocol Buffers and gRPC for high-performance internal calls.
- **Database**: SQLite with folder-based [Goose](https://github.com/pressly/goose) migrations, separating DDL and DML layers.
- **Deployment**: Multi-stage Dockerfiles and Docker Compose setups.

---

## Repository Layout

```
├── .agents/             # Agent rule configurations
├── api/
│   └── generated/       # Auto-generated OpenAPI REST controllers
├── cmd/
│   ├── migrate/         # Database migration CLI tool
│   └── server/          # Main REST & gRPC server entry point
├── internal/
│   ├── adapter/         # Primary (HTTP, gRPC) and Secondary (SQLite, CoinPaprika) adapters
│   ├── domain/          # Core domain models
│   ├── port/            # Primary (incoming) and Secondary (outgoing) interfaces
│   └── service/         # Business logic layer
├── migrations/          # goose migrations
│   ├── ddl/             # Schema definitions
│   └── dml/             # Seed data
├── proto/               # Protobuf specifications & generated Go code
├── spec/                # TypeSpec API specification files
└── smoke_test.go        # Integration smoke test suite
```

---

## Configuration

The application is configured using a combination of environment variables, a `.env` file, and command line flags.

### Environment (.env)

Create a `.env` file in the root of the project to configure external integrations:

```env
COINPAPRIKA_URL=https://api.coinpaprika.com
```

### CLI Flags

The server binary supports the following options:

- `-port`: REST HTTP port (default `9000`).
- `-grpc-port`: gRPC server port (default `9001`).
- `-db-path`: Path to SQLite database file (default `data.db`).
- `-jwt-secret`: Key used for signing and verifying JWT tokens (default `default-jwt-secret`).

---

## Usage Instructions

### Running Locally via CLI

#### 1. Compile Specifications & Build Code
Compile the TypeSpec contract and generate the corresponding OpenAPI and Protobuf Go bindings:
```bash
make compile
```

#### 2. Run Database Migrations
Apply DDL schemas and DML seeds using the custom migrations CLI:
```bash
# Run migrations up
go run ./cmd/migrate up

# Check migrations status
go run ./cmd/migrate status
```

#### 3. Start the Server
Start the REST and gRPC servers:
```bash
go run ./cmd/server --port=9000 --grpc-port=9001 --db-path=data.db --jwt-secret=my-secret-key
```

---

### Running via Docker Compose

Docker Compose builds the images and spins up both the server and migrator services. It maps the SQLite database to a persistent Docker volume (`sqlite-data`).

#### 1. Start Services
```bash
docker compose up -d
```
This automatically runs the migrator container (which applies folder-based migrations in order and exits) and spins up the server container.

#### 2. Rebuilding Containers
If you modify code or migrations, rebuild the containers using:
```bash
docker compose up -d --build
```
Or force a clean cached rebuild:
```bash
docker compose build --no-cache && docker compose up -d --force-recreate
```

#### 3. Stop Services
```bash
docker compose down
```

---

## Running Tests

### Unit & Integration Tests
Run all unit and integration tests (which utilize in-memory SQLite and hermetic HTTP mock servers):
```bash
go test ./...
```

### Smoke Test Suite
Run the dedicated integration smoke tests:
```bash
make smoke-test
```

---

## API Documentation

### REST Routes

| Method | Route | Description | Auth / Headers |
| :--- | :--- | :--- | :--- |
| **POST** | `/api/v1/login` | Login and obtain JWT token | JSON body: `{"username": "admin", "password": "password"}` |
| **POST** | `/api/v1/documents` | Create a new document | Bearer JWT Token required |
| **GET** | `/api/v1/documents/:id` | Read a document | Requires valid `PAYMENT-SIGNATURE` header (Otherwise, returns HTTP 402) |
| **GET** | `/api/v1/price/eth` | Fetch current Ethereum price | None |
| **GET** | `/api/v1/health` | Check service and database health | None |
