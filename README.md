# Golang Clean Architecture Example

**Author: Ahmet Eren BOYACI**

This is an example of a Golang project that uses Clean Architecture principles.

**Features:**
- Containerized with Docker & Docker-Compose
- Authentication with JWT
- Dependency Injection
- Centralized Error Handling
- Mocking
- Functional Tests with Fixtures (**/controller_test.go)
- Unit Tests (**/service_test.go)
- Middleware Tests (pkg/middlewares/*_test.go)
- Logging with Zap
- Caching with Redis
- Design Patterns

**Responsibility separation is applied in files as follows:**
- router.go: registering routers
- controller.go: binding parameters, calling the related service functions and returning responses
- service.go: business logic
- repository.go: database queries

## How to Run?

```bash
docker compose up -d
docker compose logs -f app
```

## How to Run Functional Tests?

```bash
docker compose up test_db -d
export MODE=test
export DB_URL=postgres://postgres:docker@localhost:5433/aeboyaci_test
go test -v {{ TEST_DIRECTORY_RELATIVE_PATH }}
```

## How to Run Unit Tests?

```bash
export MODE=test
go test -v {{ TEST_DIRECTORY_RELATIVE_PATH }}
```
