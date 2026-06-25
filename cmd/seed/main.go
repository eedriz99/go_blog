package main

import (
	"log"
	"time"

	"github.com/eedriz99/go_blog/internal/db"
	"github.com/eedriz99/go_blog/internal/env"
	"github.com/eedriz99/go_blog/internal/store"
)

func main() {
	conn, err := db.New(
		env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/go_blog?sslmode=disable"),
		env.GetInt("DB_MAX_IDLE_CONNS", 30),
		env.GetInt("DB_MAX_OPEN_CONNS", 30),
		env.GetDuration("DB_MAX_IDLE_TIME", 30*time.Second),
	)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer conn.Close()

	s := store.NewStore(conn)

	if err := db.Seed(s); err != nil {
		log.Fatalf("seeding failed: %v", err)
	}

	log.Println("seeding complete")
}
