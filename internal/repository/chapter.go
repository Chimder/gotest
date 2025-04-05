package repository

import (
	"context"

	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChapterRepository interface {
	GetChapterByMangaNameAndNumber(ctx context.Context, name string, chapter int) (ChapterRepo, error)
	ListChaptersByMangaName(ctx context.Context, animename string) ([]ChapterRepo, error)
}

type chapterRepository struct {
	db *pgxpool.Pool
}

func NewChapterRepository(db *pgxpool.Pool) ChapterRepository {
	return &chapterRepository{
		db: db,
	}
}

type ChapterRepo struct {
	Chapter   int32            `db:"chapter" json:"chapter"`
	Img       []string         `db:"img" json:"img"`
	Name      string           `db:"name" json:"name"`
	AnimeName string           `db:"animeName" json:"animeName"`
	CreatedAt pgtype.Timestamp `db:"createdAt" json:"createdAt"`
}

type GetChapterByAnimeNameAndNumberParams struct {
	AnimeName string `db:"animeName" json:"animeName"`
	Chapter   int32  `db:"chapter" json:"chapter"`
}

func (q *chapterRepository) GetChapterByMangaNameAndNumber(ctx context.Context, name string, chapter int) (ChapterRepo, error) {
	query := `SELECT * FROM "Chapter" WHERE "animeName" = @animeName AND chapter = @chapter`
	rows, err := q.db.Query(ctx, query, pgx.NamedArgs{"animeName": name, "chapter": chapter})
	if err != nil {
		return ChapterRepo{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[ChapterRepo])
}

func (q *chapterRepository) ListChaptersByMangaName(ctx context.Context, animename string) ([]ChapterRepo, error) {
	query := `SELECT * FROM "Chapter" WHERE "animeName" = $1`
	rows, err := q.db.Query(ctx, query, animename)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[ChapterRepo])
}
