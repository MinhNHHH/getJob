package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

func (app Application) ConnectReids(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
}

func (app Application) ConnectDB(dbURL string) (*sql.DB, error) {
	// Open the database connection
	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}
	// Test the database connection
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	return conn, nil
}

func ensureMigrationsDir() error {
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		// Create the "migrations" folder with permission 0755 (owner: read/write/execute, others: read/execute)
		if err := os.MkdirAll("migrations", 0755); err != nil {
			return fmt.Errorf("failed to create migrations directory: %v", err)
		}
	}
	return nil
}

func (app Application) GenerateMigration(name string) {
	err := ensureMigrationsDir()
	if err != nil {
		log.Fatal(err)
	}

	timestamp := time.Now().Format("20060102150405")
	upPath := filepath.Join("migrations", fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	downPath := filepath.Join("migrations", fmt.Sprintf("%s_%s.down.sql", timestamp, name))

	// Create up migration
	if err := os.WriteFile(upPath, []byte("-- Add your up migration here\n"), 0644); err != nil {
		log.Fatal(err)
	}

	// Create down migration
	if err := os.WriteFile(downPath, []byte("-- Add your down migration here\n"), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created migration files:\n%s\n%s\n", upPath, downPath)
}

func (app Application) Migrate(step int, uri string) {
	m, err := migrate.New(
		"file://migrations/",
		uri,
	)
	if err != nil {
		log.Fatal(err)
	}
	if step == 0 {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else {
		if err := m.Steps(step); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
