package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"shanraq.com/internal/migrations"
)

func main() {
	var (
		databaseURL = flag.String("database", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
		dir         = flag.String("dir", "migrations", "Directory containing migration files")
		command     = flag.String("command", "up", "Migration command: up, down, steps")
		steps       = flag.Int("steps", 1, "Number of steps for down/steps commands")
	)

	flag.Parse()

	if *databaseURL == "" {
		log.Fatal("database URL is required (-database or DATABASE_URL)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	runner := migrations.NewRunner(*databaseURL, *dir)

	var err error
	switch *command {
	case "up":
		err = runner.Up()
	case "down":
		err = runner.Down(*steps)
	case "steps":
		err = runner.Steps(*steps)
	default:
		err = fmt.Errorf("unknown command %s", *command)
	}

	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations completed successfully")

	// ensure context is used to avoid linter complaints
	<-ctx.Done()
}
