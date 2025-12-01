# Skeleton Application

A production-ready skeleton application for the DG Framework.

## Features

- ✅ DG Framework integration
- ✅ Database migrations with go-migrate
- ✅ Environment configuration
- ✅ HTTP routing
- ✅ Logging
- ✅ Error handling

## Quick Start

### 1. Install Dependencies

```bash
make deps
make install-migrate
```

### 2. Configure Database

Edit `config/database.yaml`:

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  database: myapp
  username: postgres
  password: secret
  schema: public  # Optional: PostgreSQL schema
```

### 3. Setup Database

```bash
# Create database
createdb myapp

# Run migrations
make migrate-up
```

### 4. Run Application

```bash
make run
```

Visit: http://localhost:8080

## Development

### Available Commands

```bash
make help              # Show all available commands
make run               # Run the application
make build             # Build the application
make test              # Run tests
make migrate-up        # Run all pending migrations
make migrate-down      # Rollback last migration
make migrate-status    # Show migration status
make migrate-create    # Create new migration
make clean             # Clean build artifacts
```

### Database Migrations

#### Run Migrations

```bash
# Run all pending migrations
make migrate-up

# Or using the CLI directly
go run cmd/migrate/main.go -direction=up
```

#### Rollback Migrations

```bash
# Rollback last migration
make migrate-down

# Or rollback multiple steps
go run cmd/migrate/main.go -direction=down -steps=2
```

#### Create New Migration

```bash
# Using make
make migrate-create NAME=add_posts_table

# Or using go-migrate CLI
migrate create -ext sql -dir database/migrations -seq add_posts_table
```

This creates:
- `database/migrations/000002_add_posts_table.up.sql`
- `database/migrations/000002_add_posts_table.down.sql`

#### Check Migration Status

```bash
make migrate-status
```

### Project Structure

```
skeleton/
├── cmd/
│   └── migrate/          # Migration CLI tool
│       └── main.go
├── config/               # Configuration files
├── database/
│   └── migrations/       # SQL migration files
│       ├── 000001_create_users_table.up.sql
│       ├── 000001_create_users_table.down.sql
│       └── README.md
├── routes/               # HTTP routes
│   ├── web.go
│   └── web_test.go
├── .env                  # Environment variables (create this)
├── go.mod
├── main.go
└── Makefile
```

## Database

### Supported Databases

- PostgreSQL (recommended)
- MySQL
- SQLite

### Migration Files

Migrations are stored in `database/migrations/` with the following format:

```
{version}_{description}.up.sql    # Forward migration
{version}_{description}.down.sql  # Rollback migration
```

Example:
```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```sql
-- 000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

### Migration Best Practices

1. **Always create both up and down migrations**
2. **Keep migrations small and focused**
3. **Test migrations before deploying**
4. **Never modify existing migrations** (create new ones instead)
5. **Use transactions when possible**

See `database/migrations/README.md` for detailed migration guide.

## Configuration

Configuration is managed through environment variables and the `.env` file.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | MyApp |
| `APP_ENV` | Environment (development, production) | development |
| `APP_PORT` | HTTP port | 8080 |
| `DB_DRIVER` | Database driver (postgres, mysql, sqlite) | postgres |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_DATABASE` | Database name | myapp |
| `DB_USERNAME` | Database username | postgres |
| `DB_PASSWORD` | Database password | - |

## Deployment

### Build

```bash
make build
```

This creates `bin/app` executable.

### Run in Production

```bash
# Set environment
export APP_ENV=production

# Run migrations
./bin/app migrate -direction=up

# Start application
./bin/app
```

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/database/migrations ./database/migrations

# Run migrations and start app
CMD ["sh", "-c", "./app migrate -direction=up && ./app"]
```

## Testing

```bash
# Run all tests
make test

# Run specific test
go test ./routes -v

# Run with coverage
go test ./... -cover
```

## CI/CD

### GitHub Actions Example

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: myapp_test
          POSTGRES_PASSWORD: secret
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run migrations
        run: go run cmd/migrate/main.go -direction=up
        env:
          DB_HOST: localhost
          DB_DATABASE: myapp_test
          DB_USERNAME: postgres
          DB_PASSWORD: secret
      
      - name: Run tests
        run: go test ./... -v
```

## Troubleshooting

### Database Connection Failed

Check your `.env` file and ensure:
- Database server is running
- Credentials are correct
- Database exists

### Migration Failed

1. Check the error message
2. Fix the SQL in the migration file
3. If partially applied, manually revert changes
4. Run migration again

### Dirty Database State

If a migration fails mid-way:

```bash
# Check status
make migrate-status

# Fix manually, then force version
migrate -path database/migrations -database $DATABASE_URL force {version}
```

## Additional Resources

- [DG Framework Documentation](https://github.com/donnigundala/dgcore)
- [go-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## License

MIT License
