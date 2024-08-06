// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chapter.sql

package queries

import (
	"context"
)

const getChapterByAnimeNameAndNumber = `-- name: GetChapterByAnimeNameAndNumber :one
SELECT chapter, img, name, "animeName", "createdAt" FROM "Chapter" WHERE "animeName" = $1 AND chapter = $2
`

type GetChapterByAnimeNameAndNumberParams struct {
	AnimeName string `db:"animeName" json:"animeName"`
	Chapter   int32  `db:"chapter" json:"chapter"`
}

func (q *Queries) GetChapterByAnimeNameAndNumber(ctx context.Context, arg GetChapterByAnimeNameAndNumberParams) (Chapter, error) {
	row := q.db.QueryRow(ctx, getChapterByAnimeNameAndNumber, arg.AnimeName, arg.Chapter)
	var i Chapter
	err := row.Scan(
		&i.Chapter,
		&i.Img,
		&i.Name,
		&i.AnimeName,
		&i.CreatedAt,
	)
	return i, err
}

const listChaptersByAnimeName = `-- name: ListChaptersByAnimeName :many
SELECT chapter, img, name, "animeName", "createdAt" FROM "Chapter" WHERE "animeName" = $1
`

func (q *Queries) ListChaptersByAnimeName(ctx context.Context, animename string) ([]Chapter, error) {
	rows, err := q.db.Query(ctx, listChaptersByAnimeName, animename)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chapter
	for rows.Next() {
		var i Chapter
		if err := rows.Scan(
			&i.Chapter,
			&i.Img,
			&i.Name,
			&i.AnimeName,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
