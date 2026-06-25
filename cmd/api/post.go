package main

import (
	"database/sql"
	"errors"
	"net/http"

	payload "github.com/eedriz99/go_blog/internal/dto/payload"
	response "github.com/eedriz99/go_blog/internal/dto/response"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
)

type postKey string

const postContextKey postKey = "post"

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	postCtx := getpostFromContext(r)
	res := response.NewPostResponse(postCtx)

	if err := writeJson(w, http.StatusOK, res); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload payload.CreatePostPayload

	if err := readJson(w, r, &payload); err != nil {
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
	postCtx := getpostFromContext(r)

	if err := readJson(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	payload.ID = postCtx.ID

	post, err := app.store.Posts.Update(ctx, payload)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	writeJson(w, http.StatusOK, response.NewPostResponse(post))
}

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

	writeJson(w, http.StatusOK, []any{})
}

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
	writeJson(w, http.StatusOK, res)
}
