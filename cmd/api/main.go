package main

import (
	"log"
	"time"

	"github.com/eedriz99/go_blog/internal/db"
	"github.com/eedriz99/go_blog/internal/env"
	"github.com/eedriz99/go_blog/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("PORT", ":8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/go_blog?sslmode=disable"),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 30*time.Second),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxIdleConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	s := store.NewStore(db)

	app := &application{
		config: cfg,
		store:  s,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
