-- name: GetMangaByName :one
SELECT * FROM "Anime" WHERE name = $1;

-- name: ListMangas :many
SELECT * FROM "Anime";

-- name: ListPopularMangas :many
SELECT * FROM "Anime" ORDER BY "ratingCount" DESC LIMIT 14;

-- name: UpdateAnimePopularity :exec
UPDATE "Anime" SET "popularity" = "popularity" + 1 WHERE "name" = $1;