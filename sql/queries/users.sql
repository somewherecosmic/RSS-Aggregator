-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: FindUserByName :one
SELECT * FROM users
WHERE name = $1;

-- name: ClearUserTable :exec
DELETE FROM users;

-- name: GetPostsForUser :many
SELECT posts.* FROM users
INNER JOIN feeds
ON feeds.user_id = users.id
INNER JOIN posts
ON posts.feed_id = feeds.id
WHERE users.id = $1
LIMIT $2;