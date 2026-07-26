package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/eedriz99/go_blog/docs"
	"github.com/eedriz99/go_blog/internal/db"
	"github.com/eedriz99/go_blog/internal/env"
	"github.com/eedriz99/go_blog/internal/store"
)

const version = "0.0.1"

var BASE_URL = fmt.Sprintf("[::]:%s/v1", env.GetString("PORT", "8000"))

// @title           Go Blog API
// @version         1.0
// @description     A blog REST API with posts, comments, and users.
// @contact.name    eedriz99
// @host            localhost:8000
// @BasePath        /v1
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
		mail: mailConfig{
			expiry: time.Hour * 10, // 10 hours activation window
		},
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
