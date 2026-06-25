# Handlers & Routes

## Registered Routes

### System

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `GET` | `/swagger/*` | `httpSwagger.Handler` | — | — |
| `GET` | `/v1/health` | `healthCheckHandler` | `health.go` | — |

### Posts

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `POST` | `/v1/posts` | `createPostHandler` | `post.go` | — |
| `GET` | `/v1/posts` | `getAllPostsHandler` | `post.go` | — |
| `GET` | `/v1/posts/{postID}` | `getPostWithCommentsHandler` | `post.go` | `postContextMiddleware` |
| `PATCH` | `/v1/posts/{postID}` | `updatePostHandler` | `post.go` | `postContextMiddleware` |
| `DELETE` | `/v1/posts/{postID}` | `deletePostHandler` | `post.go` | `postContextMiddleware` |

### Comments

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `POST` | `/v1/posts/{postID}/comments` | `createCommentHandler` | `comment.go` | `postContextMiddleware` |
| `GET` | `/v1/posts/{postID}/comments` | `getCommentsByPostHandler` | `comment.go` | `postContextMiddleware` |
| `PUT` | `/v1/comments/{commentID}` | `updateCommentHandler` | `comment.go` | — |
| `DELETE` | `/v1/comments/{commentID}` | `deleteCommentHandler` | `comment.go` | — |

### Users

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `GET` | `/v1/users/{userID}` | `GetUserHandler` | `user.go` | `userContextMiddleware` |

---

## Unregistered Handlers

These handlers exist in code but have no route wired up in `api.go`.

| Handler | File | Notes |
|---------|------|-------|
| `getPostsHandler` | `post.go:26` | Returns a single post from context — likely superseded by `getPostWithCommentsHandler` |
| `CreateUserHandler` | `user.go:13` | Empty stub — user registration not yet implemented |

---

## Middleware

| Middleware | File | Applied To |
|------------|------|------------|
| `postContextMiddleware` | `middleware.go` | All `/v1/posts/{postID}` routes |
| `userContextMiddleware` | `middleware.go` | All `/v1/users/{userID}` routes |

---

## Error Helpers

Not handlers — shared response utilities called from within handlers.

| Function | File |
|----------|------|
| `InternalServerError` | `errors.go` |
| `BadRequestError` | `errors.go` |
| `ResourceNotFoundError` | `errors.go` |
