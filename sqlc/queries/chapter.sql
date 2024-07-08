-- name: GetChapterByAnimeNameAndNumber :one
SELECT * FROM "Chapter" WHERE "animeName" = $1 AND chapter = $2;

-- name: ListChaptersByAnimeName :many
SELECT * FROM "Chapter" WHERE "animeName" = $1;
