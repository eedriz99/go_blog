package main

import (
	"log"
	"net/http"

	payload "github.com/eedriz99/go_blog/internal/dto/payload"
	response "github.com/eedriz99/go_blog/internal/dto/response"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/go-chi/chi/v5"
)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "post_id")

	var payload payload.CreateCommentPayload
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

	res := response.NewCommentWithoutUsernameResponse(comment)

	if err := writeJson(w, http.StatusCreated, res); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

func (app *application) getCommentsByPostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "post_id")
	ctx := r.Context()

	comments, err := app.store.Comments.GetByPost(ctx, postId)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	res := response.NewCommentListResponse(comments)
	writeJson(w, http.StatusOK, res)
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload payload.UpdateCommentPayload
	if err := readJson(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	payload.ID = chi.URLParam(r, "comment_id")
	payload.UserID = "cdf8c7d8-913c-4300-abee-b2165c541176" // place holder. Get it from ctx

	ctx := r.Context()
	updatedComment, err := app.store.Comments.Update(ctx, payload)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	res := response.NewCommentWithoutUsernameResponse(updatedComment)
	writeJson(w, http.StatusOK, res)
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload payload.DeleteCommentPayload
	ctx := r.Context()

	payload.ID = chi.URLParam(r, "comment_id")
	payload.UserID = "cdf8c7d8-913c-4300-abee-b2165c541176" // place holder. Get it from ctx
	log.Printf("Delete Payload: %v", payload)

	if err := app.store.Comments.Delete(ctx, payload); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	writeJson(w, http.StatusNoContent, nil)
}
