package main

import (
	"log"
	"net/http"

	"github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/dto/response"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/go-chi/chi/v5"
)

// @Summary     Create comment
// @Description Create a comment on a post
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       postID  path   string                      true "Post ID"
// @Param       payload body   payload.CreateCommentPayload    true "Comment payload"
// @Success     201 {object} response.CommentResponse
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts/{postID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postID")

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

// @Summary     Read comments for a post
// @Description Read all comments for a given post
// @Tags        comments
// @Produce     json
// @Param       postID path string true "Post ID"
// @Success     200 {array}  response.CommentResponse
// @Failure     500 {object} map[string]string
// @Router      /posts/{postID}/comments [get]
func (app *application) getCommentsByPostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postID")
	ctx := r.Context()

	comments, err := app.store.Comments.GetByPostID(ctx, postId)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	res := response.NewCommentListResponse(comments)
	writeJson(w, http.StatusOK, res)
}

// @Summary     Update comment
// @Description Update the content of a comment
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       commentID path   string                      true "Comment ID"
// @Param       payload   body   payload.UpdateCommentPayload    true "Update payload"
// @Success     200 {object} response.CommentResponse
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /comments/{commentID} [put]
func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload payload.UpdateCommentPayload
	if err := readJson(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	payload.ID = chi.URLParam(r, "commentID")
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

// @Summary     Delete comment
// @Description Delete a comment by ID
// @Tags        comments
// @Produce     json
// @Param       commentID path string true "Comment ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /comments/{commentID} [delete]
func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload payload.DeleteCommentPayload
	ctx := r.Context()

	payload.ID = chi.URLParam(r, "commentID")
	payload.UserID = "cdf8c7d8-913c-4300-abee-b2165c541176" // place holder. Get it from ctx
	log.Printf("Delete Payload: %v", payload)

	if err := app.store.Comments.Delete(ctx, payload); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	writeJson(w, http.StatusNoContent, nil)
}
