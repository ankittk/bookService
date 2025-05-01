# 📁 Project Structure: Book Search Service

This document outlines the recommended directory structure for the `booksearch` microservice and the responsibilities of each folder.

---

## 📦 Directory Layout

```txt
booksearch/
├── api/                # Versioned API definitions (optional alternative to proto/)
│   └── v1/             # Organized by version
├── cmd/                # Main applications for this project (one per binary)
│   └── server/         # The entry point for the gRPC server
│       └── main.go
├── config/             # Configuration loading logic (env, flags, etc.)
├── internal/           # Private application and business logic
│   ├── server/         # gRPC server setup and handlers
│   ├── search/         # Core domain logic for interacting with Open Library
│   └── limiter/        # Rate limiting logic
├── proto/              # .proto files and generated Go code (can split if preferred)
│   ├── v1/             # Versioned proto definitions
│   └── gen/            # Output directory for generated gRPC code
├── test/               # Integration tests (if needed)
├── scripts/            # Dev scripts, migrations, proto generation, etc.
├── Dockerfile          # Containerization
├── .env                # Optional configuration file
├── go.mod
└── go.sum
```

---

## 📂 Directory Details

### `cmd/server/`
- Entry point for the gRPC server.
- Initialize config and dependencies, start the server.
- Keep minimal; delegate logic to `internal/`.

### `config/`
- Loads and validates environment variables or `.env` files.
- Can use libraries like `viper`, `envconfig`, or `godotenv`.

### `proto/`
- Stores `.proto` definitions for the gRPC API.
- Recommended: `proto/v1/booksearch.proto` for versioning.
- Generated Go code can go in `proto/gen/`, `internal/pb/`, or `pkg/pb/`.

### `internal/`
This is the core application logic and is **not importable** by other modules.

- `internal/server/`: Implements the gRPC service, depends on `search/` and `limiter/`.
- `internal/search/`: Handles Open Library API requests and response transformation.
- `internal/limiter/`: Implements request throttling using `golang.org/x/time/rate`.

### `test/`
- Place integration or full-system tests here.
- Keep test dependencies isolated from application code.

### `scripts/`
- Contains helper shell scripts like `generate_proto.sh`, `run_dev_server.sh`, etc.

---

## 🧪 Optional/Alternative Directories

### `pkg/`
- Reusable components or libraries meant to be used across multiple services.
- Example: shared logger, middlewares.

### `build/`
- Build-specific files like Docker configs, CI/CD, Kubernetes manifests.

---

## 🧬 `internal/` vs `pkg/` in Go

| Directory    | Purpose                                                              |
|--------------|----------------------------------------------------------------------|
| `internal/`  | Contains application-specific, **private** logic. Enforced by Go.    |
| `pkg/`       | Contains **shared**, **public** packages for reuse across projects.  |

Following this structure ensures modularity, clarity, and Go best practices for microservices.
