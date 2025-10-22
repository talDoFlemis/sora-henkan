package main

import (
	"database/sql"
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/taldoflemis/sora-henkan/settings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type MigrateSettings struct {
	Database settings.DatabaseSettings `mapstructure:"database" validate:"required"`
}

func main() {
	var (
		direction string
		steps     int
	)

	flag.StringVar(&direction, "direction", "up", "Migration direction: up, down, or force")
	flag.IntVar(&steps, "steps", 0, "Number of steps to migrate (0 = all)")
	flag.Parse()

	// Load settings
	cfg, err := settings.LoadConfig[MigrateSettings]("MIGRATE", settings.BaseSettings)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Build connection string from settings
	dbURL := cfg.Database.BuildConnectionString()

	if err := runMigrations(dbURL, direction, steps); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

func runMigrations(dbURL, direction string, steps int) error {
	// Open database connection
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create source from embedded FS
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if dirty {
		log.Printf("WARNING: Database is in dirty state (version: %d)", version)
	} else if !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("Current migration version: %d", version)
	} else {
		log.Println("No migrations applied yet")
	}

	// Run migration based on direction
	switch direction {
	case "up":
		if steps > 0 {
			log.Printf("Migrating up %d steps...", steps)
			if err := m.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("migration up failed: %w", err)
			}
		} else {
			log.Println("Migrating up to latest version...")
			if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("migration up failed: %w", err)
			}
		}

	case "down":
		if steps > 0 {
			log.Printf("Migrating down %d steps...", steps)
			if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("migration down failed: %w", err)
			}
		} else {
			log.Println("Migrating down all...")
			if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("migration down failed: %w", err)
			}
		}

	case "force":
		if steps == 0 {
			return fmt.Errorf("must specify version with -steps flag when using force")
		}
		log.Printf("Forcing version to %d...", steps)
		if err := m.Force(steps); err != nil {
			return fmt.Errorf("force version failed: %w", err)
		}

	default:
		return fmt.Errorf("invalid direction: %s (use up, down, or force)", direction)
	}

	// Get new version
	newVersion, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get new version: %w", err)
	}

	if dirty {
		log.Printf("WARNING: Database is now in dirty state (version: %d)", newVersion)
	} else if !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("New migration version: %d", newVersion)
	}

	return nil
}
