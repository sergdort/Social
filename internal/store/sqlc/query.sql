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