package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eedriz99/go_blog/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Get("/", app.getAllPostsHandler)
			r.Route("/{post_id}", func(r chi.Router) {
				r.Get("/", app.getPostsHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
					r.Get("/", app.getCommentsByPostHandler)

				})
			})
		})
	})
	return r
}

func (app *application) run(router http.Handler) error {
	//mux := router()
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute * 1,
	}
	log.Printf("Listening on port %s\n", app.config.addr)
	return srv.ListenAndServe()
}
