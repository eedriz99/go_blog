package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/dto/response"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
)

type postKey string

const postContextKey postKey = "post"

// @Summary     Read posts by user
// @Description Read all posts belonging to a specific user
// @Tags        users
// @Produce     json
// @Param       userID path string true "User ID"
// @Success     200 {object} response.PostListResponse
// @Failure     500 {object} map[string]string
// @Router      /users/{userID}/posts [get]
func (app *application) getPostsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserFromContext(r).ID
	posts, err := app.store.Posts.GetByUserID(r.Context(), userID)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, response.NewListPostResponse(posts)); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

}

// @Summary     Create post
// @Description Create a new blog post
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param       payload body payload.CreatePostPayload true "Post payload"
// @Success     201 {object} response.PostResponse
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload payload.CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
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
	if err := writeJSON(w, http.StatusCreated, res); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

// @Summary     Read all posts
// @Description Read all blog posts
// @Tags        posts
// @Produce     json
// @Success     200 {array}  response.PostResponse
// @Failure     500 {object} map[string]string
// @Router      /posts [get]
func (app *application) getAllPostsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	posts, err := app.store.Posts.GetAll(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusOK, response.NewListPostResponse([]model.Post{}))
			return
		}
		app.InternalServerError(w, r, err)
		return
	}

	response := response.NewListPostResponse(posts)

	writeJSON(w, http.StatusOK, response)
}

// @Summary     Update post
// @Description Update title, content, or tags of an existing post
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param       postID  path   string                    true "Post ID"
// @Param       payload body   payload.UpdatePostPayload     true "Update payload"
// @Success     200 {object} response.PostResponse
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload payload.UpdatePostPayload
	ctx := r.Context()
	postCtx := getpostFromContext(r)

	if err := readJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	payload.ID = postCtx.ID

	post, err := app.store.Posts.Update(ctx, payload)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, response.NewPostResponse(post))
}

// @Summary     Delete post
// @Description Delete a post by ID
// @Tags        posts
// @Produce     json
// @Param       postID path string true "Post ID"
// @Success     200 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	postCtx := getpostFromContext(r)

	if err := app.store.Posts.Delete(ctx, postCtx.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			// Race condition can cause the post to be deleted after the
			// postContextMiddleware retrieves it and before this handler executes,
			// so we should return 404 if the post is not found during deletion
			app.ResourceNotFoundError(w, r, err)
		default:
			app.InternalServerError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, []any{})
}

// @Summary     Read post with comments
// @Description Read a post and all its comments
// @Tags        posts
// @Produce     json
// @Param       postID path string true "Post ID"
// @Success     200 {object} response.PostWithCommentsResponse
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts/{postID} [get]
func (app *application) getPostWithCommentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getpostFromContext(r)

	comments, err := app.store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			comments = []store.CommentWithUsername{}
		default:
			app.InternalServerError(w, r, err)
			return
		}
	}

	res := response.NewPostWithCommentsResponse(response.NewPostResponse(post), comments)
	writeJSON(w, http.StatusOK, res)
}
