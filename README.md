# 📚 Book Search Service

A simple microservice that allows searching for books using the [Open Library API](https://openlibrary.org/developers/api).
The service supports search by title, author, and other filters, along with built-in pagination and rate limiting.

---

## 🚀 Features

### 🔍 gRPC Endpoint
- `SearchBooks` RPC method that allows:
	- Search by **title**, **author**, and optional filters like **language** or **publication year**.

### 📄 Pagination
- Support for paginated results using `page` and `limit` parameters.
- Helps reduce response size and improve performance.

### 🚦 Rate Limiting
- Prevents API abuse with rate limits per client IP or API key.
- Example: **5 requests per minute**.
- Uses in-memory rate limiter (`golang.org/x/time/rate`).

### 🧯 Error Handling
- Maps common issues to appropriate gRPC status codes:
	- `INVALID_ARGUMENT` – for bad parameters.
	- `NOT_FOUND` – if no matching results are found.
	- `INTERNAL` – for unexpected errors.

### 📦 API Structure
- Well-defined `.proto` file.
- Rich response structure including:
	- Total results count.
	- Current page number.
	- Total pages.
	- Books list.

### ⚙️ Configuration
- Load settings via environment variables:
	- `PORT` (default: 8080)
	- `LOG_LEVEL` (e.g., INFO, DEBUG)
- Optionally supports a config file for advanced settings (e.g., rate limits).

### ✅ Unit Testing
- Unit tests for:
	- Search logic.
	- Rate limiting logic.
- External API calls are mocked to ensure fast and reliable test execution.
- Integration tests for gRPC server.
- Mock server for testing client interactions.
- Test coverage report generation.
- Code quality checks using `golangci-lint`.
- Code formatting with `go fmt`.
- Static analysis with `staticcheck`.
- Dependency management with `go mod tidy`.
- Security checks with `gosec`.
- Performance profiling with `pprof`.
- Benchmarking for critical functions.

---

## 📁 Project Structure

Refer to the [docs/](./docs/) folder for detailed documentation of each step and design decision, including:

- [Package Structure](./docs/PackageStructure.md)
- [Design Principles](./docs/DesignPrinciples.md)
- [API Design](./docs/APIDesign.md)

---

## 🧪 Running & Testing

### Prerequisites
- Go ≥ 1.20
- `protoc` compiler with gRPC Go plugin

### Running the Service
```bash
export PORT=8080
export LOG_LEVEL=INFO
go run cmd/server/main.go
