# Handlers & Routes

Endpoints organised by HTTP method: **GET → POST → PATCH/PUT → DELETE**

---

## GET — Read

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `GET` | `/v1/health` | `healthCheckHandler` | `health.go` | — |
| `GET` | `/v1/posts` | `getAllPostsHandler` | `post.go` | — |
| `GET` | `/v1/posts/{postID}` | `getPostWithCommentsHandler` | `post.go` | `postContextMiddleware` |
| `GET` | `/v1/posts/{postID}/comments` | `getCommentsByPostHandler` | `comment.go` | `postContextMiddleware` |
| `GET` | `/v1/users/{userID}` | `GetUserHandler` | `user.go` | `userContextMiddleware` |
| `GET` | `/v1/users/{userID}/posts` | `getPostsByUserIDHandler` | `post.go` | `userContextMiddleware` |
| `GET` | `/swagger/*` | `httpSwagger.Handler` | — | — |

---

## POST — Create

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `POST` | `/v1/posts` | `createPostHandler` | `post.go` | — |
| `POST` | `/v1/posts/{postID}/comments` | `createCommentHandler` | `comment.go` | `postContextMiddleware` |

---

## PATCH / PUT — Update

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `PATCH` | `/v1/posts/{postID}` | `updatePostHandler` | `post.go` | `postContextMiddleware` |
| `PUT` | `/v1/comments/{commentID}` | `updateCommentHandler` | `comment.go` | — |

---

## DELETE — Delete

| Method | Path | Handler | File | Middleware |
|--------|------|---------|------|------------|
| `DELETE` | `/v1/posts/{postID}` | `deletePostHandler` | `post.go` | `postContextMiddleware` |
| `DELETE` | `/v1/comments/{commentID}` | `deleteCommentHandler` | `comment.go` | — |

---

## Unregistered Handlers

Exist in code but not yet wired to a route.

| Handler | File | Notes |
|---------|------|-------|
| `CreateUserHandler` | `user.go` | Empty stub — user registration not yet implemented |

---

## Middleware

| Middleware | File | Applied To |
|------------|------|------------|
| `postContextMiddleware` | `middleware.go` | All `/v1/posts/{postID}` routes |
| `userContextMiddleware` | `middleware.go` | All `/v1/users/{userID}` routes |

---

## Error Helpers

Shared response utilities — not handlers.

| Function | File |
|----------|------|
| `InternalServerError` | `errors.go` |
| `BadRequestError` | `errors.go` |
| `ResourceNotFoundError` | `errors.go` |
