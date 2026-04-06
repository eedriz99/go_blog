package main

import (
	"database/sql"
	"errors"
	"net/http"

	payload "github.com/eedriz99/go_blog/internal/dto/payload"
	response "github.com/eedriz99/go_blog/internal/dto/response"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := chi.URLParam(r, "post_id")

	m, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	res := response.NewPostResponse(m)

	if err := writeJson(w, http.StatusOK, res); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload payload.CreatePostPayload

	if err := readJson(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	post := &model.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  "cdf8c7d8-913c-4300-abee-b2165c541176", // placeholder value
	}

	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	res := response.NewPostResponse(post)
	if err := writeJson(w, http.StatusCreated, res); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

func (app *application) getAllPostsHandler(w http.ResponseWriter, r *http.Request) {

	UserID := "cdf8c7d8-913c-4300-abee-b2165c541176" // placeholder value should be taken from context

	ctx := r.Context()

	posts, err := app.store.Posts.GetAll(ctx, UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, http.StatusOK, response.NewListPostResponse([]model.Post{}))
			return
		}
		app.InternalServerError(w, r, err)
		return
	}

	response := response.NewListPostResponse(posts)

	writeJson(w, http.StatusOK, response)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload payload.UpdatePostPayload
	ctx := r.Context()
	payload.ID = chi.URLParam(r, "post_id")

	if err := readJson(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	m, err := app.store.Posts.Update(ctx, payload)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	writeJson(w, http.StatusOK, m)
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	payload := payload.DeletePostPayload{
		ID:     chi.URLParam(r, "post_id"),
		UserID: "cdf8c7d8-913c-4300-abee-b2165c541176", // placeholder value
	}

	if err := app.store.Posts.Delete(ctx, payload); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.BadRequestError(w, r, err)
			return
		default:
			app.InternalServerError(w, r, err)
		}

	}

	writeJson(w, http.StatusOK, []any{})
}
