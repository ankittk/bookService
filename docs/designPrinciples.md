# SOLID Principles in Go (Golang)

This document explains how to apply SOLID principles using idiomatic Go, with examples.

---

## 1. **SRP: Single Responsibility Principle**

> Every package, file, or type should have one reason to change.

Each package should be responsible for only one thing. Avoid mixing concerns such as logging, configuration, or database logic within the same package.

**Example:**

```go
// Package db Package: db
package db

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id int) error
}
```

```go
// Package logger Package: logger
package logger

type Logger interface {
	Log(message string)
}
```

By separating responsibilities, each package can evolve independently.

---

## 2. **OCP: Open/Closed Principle**

> Code should be open for extension but closed for modification.

You should be able to add new functionality without modifying existing code. Interfaces allow for new implementations without altering existing logic.

**Example:**

```go
package auth

type Authenticator interface {
	Authenticate(username, password string) (bool, error)
}

type BasicAuthenticator struct{}

func (b BasicAuthenticator) Authenticate(username, password string) (bool, error) {
	if username == "admin" && password == "password" {
		return true, nil
	}
	return false, nil
}

type OAuthAuthenticator struct{}

func (o OAuthAuthenticator) Authenticate(username, password string) (bool, error) {
	if username == "user" && password == "oauth_token" {
		return true, nil
	}
	return false, nil
}
```

You can introduce `JWTAuthenticator` or `LDAPAuthenticator` without modifying existing authenticator logic.

---

## 3. **LSP: Liskov Substitution Principle**

> Subtypes should be replaceable for their base types without breaking functionality.

Any implementation of an interface should be substitutable without altering the correctness of the program.

**Example:**

```go
type Logger interface {
	Log(message string)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Log(message string) {
	fmt.Println(message)
}

type FileLogger struct{}

func (f FileLogger) Log(message string) {
	file, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString(message + "\n")
	writer.Flush()
}

func LogMessage(logger Logger, message string) {
	logger.Log(message)
}
```

You can pass either `ConsoleLogger` or `FileLogger` without changing `LogMessage`.

---

## 4. **ISP: Interface Segregation Principle**

> Don't force clients to depend on interfaces they do not use.

Design small, specific interfaces instead of large, all-in-one interfaces.

**Example:**

```go
type Logger interface {
	Log(message string)
}

type Database interface {
	Connect()
	Query(query string) ([]string, error)
	Close()
}
```

This allows users to implement only what they need.

---

## 5. **DIP: Dependency Inversion Principle**

> High-level modules should not depend on low-level modules. Both should depend on abstractions.

Use interfaces to abstract and inject dependencies, making your code loosely coupled and easier to test.

**Example:**

```go
type Limiter interface {
	Limit() bool
}

type Searcher interface {
	Search(query string) ([]string, error)
}

type server struct {
	limiter  Limiter
	searcher Searcher
}
```

By injecting abstractions, the `server` can work with any implementation of `Limiter` and `Searcher`.

---

## Summary

| Principle | Description                                        |
|-----------|----------------------------------------------------|
| SRP       | One reason to change per module/package            |
| OCP       | Open for extension, closed for modification        |
| LSP       | Replace base types with derived types safely       |
| ISP       | Prefer small, specific interfaces                  |
| DIP       | Rely on abstractions, not concrete implementations |


# ✅ Go Project Design Principles
This project follows idiomatic Go practices to ensure maintainability, readability, and testability. 
Below are key Go-specific design principles that guide the architecture:
---

## 1. Interfaces Are Satisfied Implicitly

- Interfaces in Go are implemented implicitly — no `implements` keyword.
- **Define interfaces where they are used**, not where they are implemented.
- Example: Define the `BookSearcher` interface in the `internal/server`, not in `internal/search`.

> 🧠 *Avoid:* defining large interfaces like `type BookService interface { ... }` at the data layer.

---

## 2. The Smaller the Interface, the Better

- Prefer small, focused interfaces.
- Example: `io.Reader`, `io.Writer` — very small and composable.
- Your interfaces (e.g., `BookSearcher`) should ideally have **one method**.

> 🧠 *Avoid:* interfaces with unused methods.

---

## 3. Don’t Overuse Interfaces

- Use **concrete types** by default.
- Use interfaces **only when you need abstraction** — typically for testing or plugability.

> 🧠 *Avoid:* defining interfaces just to enable mocking unless absolutely necessary.

---

## 4. Composition Over Inheritance

- Go favors **composition** over inheritance (which Go doesn’t have).
- Build functionality using **struct embedding** or **interface composition**.

> ✅ Example:  
> `type Server struct { searcher BookSearcher }` — composing via interface.

---

## 5. Export Only What You Need

- Use **lowercase names** for unexported symbols.
- Everything under `internal/` is already private to your module.

> 🧠 *Avoid:* exporting structs, functions, or interfaces unnecessarily.

---

## 6. Keep `main.go` Minimal

- `main.go` should only:
	- Parse config
	- Wire dependencies
	- Start the server
- All logic lives in `internal/...`.

---

## 7. Handle Errors Explicitly, Early, and Often

- Go doesn’t use exceptions. Always return and check `error`.
- Prefer **early returns** for cleaner flow.
- Add **context** to errors when needed.

---

## 8. No Fancy DI Frameworks

- Go embraces **manual dependency injection**.
- Use struct constructors or function parameters to pass dependencies.

> 🧠 Manual DI is explicit, simple, and easy to follow — no magic.

---

## 9. Build Once, Run Anywhere

- Go produces static binaries with no runtime dependencies.
- This makes it ideal for containers and lightweight deployment.

> ✅ Result: Small, secure, and fast containers.

---

## 10. Favor Simplicity Over Cleverness

- Go prioritizes **clear and boring code** over abstraction and trickiness.
- The simplest solution that works is usually the best one.

> 🧠 Readability > Cleverness
