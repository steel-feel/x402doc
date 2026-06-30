# Workspace Rules & Invariants

## Docker Configurations
- **Multi-Binary CMD Invariant:** For Dockerfiles packaging multiple binaries (e.g. server and migrator), always use `CMD` instead of `ENTRYPOINT`. This ensures that Docker Compose can override the run command.

## Database Migrations
- **DDL & DML Separation:** Separate migrations into separate files: DDL files for schema modifications and DML files for inserting seed data.

## TypeSpec Specifications
- **Route Collisions:** Explicitly decorate all HTTP operations with `@route` (e.g., `@route("login")`) to prevent route conflicts.
- **Protobuf Generation:**
  - Avoid creating models named `[OperationName]Request` as they collide with auto-generated RPC wrappers.
  - Decorate all parameters of Protobuf interface methods with `@field(index)`.
