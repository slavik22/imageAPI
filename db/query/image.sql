-- name: CreateImage :one
INSERT INTO images (
    user_id,
    image_path,
    image_url
) VALUES ($1, $2, $3) RETURNING *;

-- name: GetImages :many
SELECT * FROM images
WHERE user_id = $1;
