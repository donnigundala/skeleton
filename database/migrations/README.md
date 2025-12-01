# Database Migrations

This directory contains database migrations managed by [go-migrate](https://github.com/golang-migrate/migrate).

## Migration Files

Migrations are stored as SQL files with the following naming convention:
```
{version}_{description}.up.sql
{version}_{description}.down.sql
```

Example:
```
000001_create_users_table.up.sql
000001_create_users_table.down.sql
000002_add_posts_table.up.sql
000002_add_posts_table.down.sql
```

## Running Migrations

### Using the CLI Tool

```bash
# Run all pending migrations
go run cmd/migrate/main.go -direction=up

# Rollback last migration
go run cmd/migrate/main.go -direction=down -steps=1

# Rollback all migrations
go run cmd/migrate/main.go -direction=down

# Migrate to specific version
go run cmd/migrate/main.go -direction=up -version=2

# Check current version
go run cmd/migrate/main.go -direction=version
```

### Using Make Commands

```bash
# Run migrations
make migrate-up

# Rollback one migration
make migrate-down

# Check migration status
make migrate-status

# Create new migration
make migrate-create NAME=add_posts_table
```

## Creating New Migrations

### Using go-migrate CLI

```bash
# Install go-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration
migrate create -ext sql -dir database/migrations -seq create_posts_table
```

This creates:
- `{version}_create_posts_table.up.sql`
- `{version}_create_posts_table.down.sql`

### Manual Creation

1. Determine next version number (e.g., `000002`)
2. Create two files:
   - `000002_description.up.sql` - Forward migration
   - `000002_description.down.sql` - Rollback migration

## Best Practices

1. **Always create both up and down migrations**
   - Up: Apply the change
   - Down: Revert the change

2. **Keep migrations small and focused**
   - One logical change per migration
   - Easier to debug and rollback

3. **Test migrations**
   - Test both up and down migrations
   - Test on a copy of production data

4. **Use transactions when possible**
   ```sql
   BEGIN;
   -- Your migration here
   COMMIT;
   ```

5. **Never modify existing migrations**
   - Once applied to production, create a new migration instead

6. **Use descriptive names**
   - Good: `create_users_table`, `add_email_index`
   - Bad: `migration1`, `update`

## Migration Structure

### Up Migration Example

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

### Down Migration Example

```sql
-- 000001_create_users_table.down.sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

## Troubleshooting

### Dirty Database State

If a migration fails mid-way, the database may be in a "dirty" state:

```bash
# Check status
go run cmd/migrate/main.go -direction=version

# Fix manually, then force version
migrate -path database/migrations -database $DATABASE_URL force {version}
```

### Migration Failed

1. Check the error message
2. Fix the SQL in the migration file
3. If migration was partially applied:
   - Manually revert changes
   - Or create a new migration to fix

### Reset Database

**WARNING: This will delete all data!**

```bash
# Rollback all migrations
go run cmd/migrate/main.go -direction=down

# Or drop and recreate database
dropdb myapp
createdb myapp

# Run migrations again
go run cmd/migrate/main.go -direction=up
```

## Environment Variables

Configure database connection via environment variables:

```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=myapp
DB_USERNAME=postgres
DB_PASSWORD=secret
```

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Run Migrations
  run: go run cmd/migrate/main.go -direction=up
  env:
    DB_HOST: ${{ secrets.DB_HOST }}
    DB_DATABASE: ${{ secrets.DB_DATABASE }}
    DB_USERNAME: ${{ secrets.DB_USERNAME }}
    DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

### Docker Example

```dockerfile
# Run migrations on container start
CMD ["sh", "-c", "go run cmd/migrate/main.go -direction=up && ./app"]
```

## Additional Resources

- [go-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Migration Best Practices](https://www.postgresql.org/docs/current/ddl-alter.html)
