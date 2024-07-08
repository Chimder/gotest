-- name: GetUserByEmail :one
SELECT * FROM "User" WHERE "email" = $1;

-- name: GetUserFavoritesByEmail :one
SELECT "favorite" FROM "User" WHERE "email" = $1;

-- name: GetAnimeByNames :many
SELECT * FROM "Anime" WHERE "name" = ANY($1::text[]);

-- name: DeleteUserByEmail :exec
DELETE FROM "User" WHERE "email" = $1;

-- name: InsertUser :one
INSERT INTO "User" (id, email, name, image)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUserFavorites :exec
UPDATE "User" SET "favorite" = $1 WHERE "email" = $2;