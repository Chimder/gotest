package repository

import (
	"context"

	"github.com/chimas/GoProject/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChapterRepository interface {
	GetChapterByMangaNameAndNumber(ctx context.Context, name string, chapter int) (models.ChapterRepo, error)
	ListChaptersByMangaName(ctx context.Context, mangaName string) ([]models.ChapterRepo, error)
}

type chapterRepository struct {
	db *pgxpool.Pool
}

func NewChapterRepository(db *pgxpool.Pool) ChapterRepository {
	return &chapterRepository{
		db: db,
	}
}

func (q *chapterRepository) GetChapterByMangaNameAndNumber(ctx context.Context, name string, chapter int) (models.ChapterRepo, error) {
	query := `SELECT * FROM "Chapter" WHERE "animeName" = @animeName AND chapter = @chapter`
	rows, err := q.db.Query(ctx, query, pgx.NamedArgs{"animeName": name, "chapter": chapter})
	if err != nil {
		return models.ChapterRepo{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[models.ChapterRepo])
}

func (q *chapterRepository) ListChaptersByMangaName(ctx context.Context, mangaName string) ([]models.ChapterRepo, error) {
	query := `SELECT * FROM "Chapter" WHERE "animeName" = $1`
	rows, err := q.db.Query(ctx, query, mangaName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.ChapterRepo])
}
