package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.ResourceNotFoundError(w, r, err)

			default:
				app.InternalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, postContextKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.ResourceNotFoundError(w, r, err)

			default:
				app.InternalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// <======================= Helper functions =====================>

func getpostFromContext(r *http.Request) *model.Post {
	post, ok := r.Context().Value(postContextKey).(*model.Post)
	if !ok {
		panic("postContextMiddleware not applied to this route")
	}
	return post
}

func getUserFromContext(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("userContextMiddleware not applied to this route")
	}
	return user
}
