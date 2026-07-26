package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eedriz99/go_blog/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
	mail mailConfig
}

type mailConfig struct {
	expiry time.Duration
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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"operationsSorter": `(a, b) => {
				const order = { 'get': 1, 'post': 2, 'patch': 3, 'put': 3, 'delete': 4 };
				return (order[a.get('method')] || 5) - (order[b.get('method')] || 5);
			}`,
		}),
	))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Get("/", app.getAllPostsHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware) // add post to context for handlers that need it
				r.Get("/", app.getPostWithCommentsHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
					r.Get("/", app.getCommentsByPostHandler)
				})
			})
		})
		r.Route("/comments", func(r chi.Router) {
			r.Put("/{commentID}", app.updateCommentHandler)
			r.Delete("/{commentID}", app.deleteCommentHandler)
		})
		r.Route("/users/{userID}", func(r chi.Router) {
			r.Use(app.userContextMiddleware)
			r.Get("/", app.GetUserHandler)
			r.Get("/posts", app.getPostsByUserIDHandler)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.CreateUserHandler)
			r.Post("/login", app.LoginHandler)
			r.Post("/activate/{token}", app.ActivateUserHandler)
		})
	})
	return r
}

func (app *application) run(router http.Handler) error {
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
