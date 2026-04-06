package main

import (
	"log"
	"net/http"

	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "post_id")

	var payload store.CreateCommentPayload
	if err := readJson(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	comment := &model.Comment{
		Content: payload.Content,
		PostID:  postId,
		UserID:  "cdf8c7d8-913c-4300-abee-b2165c541176", // place holder. Get it from ctx
	}

	log.Printf("Comment: %v", comment)

	ctx := r.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusCreated, comment); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}
