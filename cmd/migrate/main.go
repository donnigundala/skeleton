package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/donnigundala/dg-core/config"
	dgdb "github.com/donnigundala/dg-database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Parse flags
	var (
		direction = flag.String("direction", "up", "Migration direction: up, down, version")
		steps     = flag.Int("steps", 0, "Number of steps to migrate (for down)")
		version   = flag.Uint("version", 0, "Migrate to specific version")
		configDir = flag.String("config", "config", "Configuration directory")
	)
	flag.Parse()

	// Load configuration
	if err := config.LoadWithPaths(*configDir); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Build database configuration
	port := 5432 // default
	if p := config.Get("database.port"); p != nil {
		if portInt, ok := p.(int); ok {
			port = portInt
		}
	}

	dbConfig := dgdb.DefaultConfig().
		WithDriver(config.GetString("database.driver")).
		WithHost(config.GetString("database.host")).
		WithPort(port).
		WithDatabase(config.GetString("database.database")).
		WithCredentials(
			config.GetString("database.username"),
			config.GetString("database.password"),
		)

	// Add schema if specified (PostgreSQL)
	if schema := config.GetString("database.schema"); schema != "" {
		dbConfig = dbConfig.WithSchema(schema)
	}

	// Create database manager
	manager, err := dgdb.NewManager(dbConfig, nil)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer manager.Close()

	// Get SQL database
	sqlDB, err := manager.SQL()
	if err != nil {
		log.Fatalf("Failed to get SQL database: %v", err)
	}

	// Create database driver based on driver type
	var driver database.Driver
	switch dbConfig.Driver {
	case "postgres":
		driver, err = postgres.WithInstance(sqlDB, &postgres.Config{})
	case "mysql":
		driver, err = mysql.WithInstance(sqlDB, &mysql.Config{})
	default:
		log.Fatalf("Unsupported database driver: %s (supported: postgres, mysql)", dbConfig.Driver)
	}

	if err != nil {
		log.Fatalf("Failed to create database driver: %v", err)
	}

	// Create migrator
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		dbConfig.Database,
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	// Run migration based on direction
	switch *direction {
	case "up":
		if *version > 0 {
			err = m.Migrate(*version)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration failed: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to run")
		} else {
			fmt.Println("Migrations completed successfully")
		}

	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Rollback failed: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
		} else {
			fmt.Println("Rollback completed successfully")
		}

	case "version":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Failed to get version: %v", err)
		}
		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations applied yet")
		} else {
			fmt.Printf("Current version: %d\n", version)
			if dirty {
				fmt.Println("WARNING: Database is in dirty state!")
			}
		}
		return

	default:
		log.Fatalf("Invalid direction: %s (use: up, down, version)", *direction)
	}

	// Show current version
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Warning: Failed to get current version: %v", err)
	} else if err != migrate.ErrNilVersion {
		fmt.Printf("Current version: %d\n", currentVersion)
		if dirty {
			fmt.Println("WARNING: Database is in dirty state!")
		}
	}
}
