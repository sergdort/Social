-- name: DeleteFollow :execrows
DELETE
FROM followers
WHERE user_id = $1
  AND follower_id = $2;

-- name: CreateFollow :exec
INSERT INTO followers (user_id, follower_id)
VALUES ($1, $2);

-- name: GetAllCommentsByPostID :many
SELECT c.id,
       c.post_id,
       c.user_id,
       c.content,
       c.created_at,
       u.username,
       u.id
FROM comments c
         JOIN users u ON u.id = c.user_id
WHERE c.post_id = $1
ORDER BY c.created_at DESC;

-- name: CreateComment :one
INSERT INTO comments (post_id, user_id, content)
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: DeleteUserInvitationByUserID :exec
DELETE
FROM user_invitations
WHERE user_id = $1;

-- name: DeleteUserByID :exec
DELETE
FROM users
WHERE id = $1;

-- name: DeleteUserInvitationByToken :exec
DELETE
FROM user_invitations
WHERE token = $1;

-- name: CreateUserInvitation :exec
INSERT INTO user_invitations (token, user_id, expiry)
VALUES ($1, $2, $3);

-- name: UpdatePost :one
UPDATE posts
SET content = $1,
    title   = $2,
    version = version + 1
WHERE id = $3
  AND version = $4
RETURNING version;

-- name: DeletePostByID :execrows
DELETE
FROM posts
WHERE id = $1;

-- name: CreatePost :one
INSERT INTO posts (content, title, user_id, tags)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetRoleByName :one
SELECT id, name, description, level
FROM roles
WHERE name = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, password, role_id)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at;

-- name: GetUserByID :one
SELECT users.id,
       users.username,
       users.email,
       users.created_at,
       users.is_active,
       r.id          as role_id,
       r.name        as role_name,
       r.description as role_description,
       r.level       as role_level
FROM users
         JOIN roles r ON (users.role_id = r.id)
WHERE users.id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password, username, created_at, is_active
FROM users
WHERE email = $1;

-- name: ActiveUserByInvitationToken :exec
UPDATE users u
SET is_active = TRUE
FROM user_invitations i
WHERE i.user_id = u.id
  AND i.token = $1
  AND i.expiry > $2;

-- name: GetPostByID :one
SELECT id,
       content,
       title,
       user_id,
       created_at,
       updated_at,
       tags,
       version
FROM posts
WHERE id = $1;

-- name: GetUserFeed :many
SELECT p.id,
       p.user_id,
       p.title,
       p.content,
       p.created_at,
       p.tags,
       COUNT(c.id) AS comments_count,
       u.username
FROM posts p
         LEFT JOIN comments c ON c.post_id = p.id
         LEFT JOIN users u ON p.user_id = u.id
         JOIN followers f ON f.follower_id = p.user_id
    OR p.user_id = $1
WHERE (f.user_id = $1 OR p.user_id = $1)
  AND ($4 = '' OR LOWER(p.title) LIKE LOWER('%' || $4 || '%') OR LOWER(p.content) LIKE LOWER('%' || $4 || '%'))
  AND (p.tags @> $5 OR $5 = '{}')
GROUP BY p.id, u.username
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;