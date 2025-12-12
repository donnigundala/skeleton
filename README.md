# Skeleton Application

A production-ready skeleton application for the DG Framework.

## Features

- ✅ **DG Framework Integration**: Core, Database, Redis, Queue, Filesystem
- ✅ **Docker Ready**: Postgres, Redis, MinIO included
- ✅ **Database Migrations**: Automated version control for schema
- ✅ **Rich CLI**: Makefile for common tasks
- ✅ **Production Grade**: Logging, Error Handling, Graceful Shutdown

## 🚀 Quick Start

Get up and running in **2 minutes**.

### 1. Configure Environment
```bash
cp .env.example .env
```

### 2. Setup
This installs dependencies, starts Docker containers (DB, Redis, MinIO), and runs migrations.
```bash
make setup
```

### 3. Run
```bash
make run
```
Visit http://localhost:8080/health

## 🛠️ Development

### Local Services (Docker)
We use `docker-compose` to run development dependencies.

- **Postgres**: localhost:5432 (user: postgres, pass: secret)
- **Redis**: localhost:6379
- **MinIO** (S3): localhost:9001 (Console), localhost:9000 (API)

Manage with:
```bash
make docker-up    # Start
make docker-down  # Stop
```

### CLI Commands (Makefile)

| Command | Description |
|---------|-------------|
| `make run` | Run the application |
| `make test` | Run all tests |
| `make migrate-up` | Run pending migrations |
| `make migrate-create NAME=x` | Create new migration |
| `make clean` | Clean build artifacts |

### Project Structure

```
skeleton/
├── app/                  # Application Logic
│   ├── http/             # Controllers, Middleware
│   │   ├── routes/       # Route definitions
│   │   └── controllers/  # Request handlers
│   ├── jobs/             # Background jobs
│   ├── models/           # Domain models
│   ├── services/         # Business logic
│   └── providers/        # Service providers
├── bootstrap/            # App initialization
├── config/               # Configuration files (yaml)
├── database/
│   └── migrations/       # SQL migrations
├── cmd/                  # CLI commands
├── docker-compose.yml    # Local dev stack
├── Makefile              # Task runner
└── main.go               # Entry point
```

## 📚 Documentation

- [DG Framework](https://github.com/donnigundala/dg-core)
- [Migrations](database/migrations/README.md)
- [API Documentation](docs/api.md) (Coming soon)

## License

MIT License
