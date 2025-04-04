package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MangaRepository interface {
	GetMangaByName(ctx context.Context, name string) (MangaRepo, error)
	GetMangaByNames(ctx context.Context, names []string) ([]MangaRepo, error)
	ListMangas(ctx context.Context) ([]MangaRepo, error)
	ListPopularMangas(ctx context.Context) ([]MangaRepo, error)
	UpdateMangaPopularity(ctx context.Context, name string) error
}

type mangaRepository struct {
	db *pgxpool.Pool
}

func NewMangaRepository(db *pgxpool.Pool) MangaRepository {
	return &mangaRepository{
		db: db,
	}
}

func (q *mangaRepository) GetMangaByName(ctx context.Context, name string) (MangaRepo, error) {
	query := `SELECT * FROM "Anime" WHERE name = $1`
	rows, err := q.db.Query(ctx, query, name)
	if err != nil {
		return MangaRepo{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[MangaRepo])
}

func (q *mangaRepository) ListMangas(ctx context.Context) ([]MangaRepo, error) {
	query := `SELECT * FROM "Anime"`
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[MangaRepo])
}

func (q *mangaRepository) GetMangaByNames(ctx context.Context, names []string) ([]MangaRepo, error) {
	query := `SELECT * FROM "Anime" WHERE "name" = ANY(@name::text[])`
	rows, err := q.db.Query(ctx, query, pgx.NamedArgs{"name": names})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[MangaRepo])
}

func (q *mangaRepository) ListPopularMangas(ctx context.Context) ([]MangaRepo, error) {
	query := `SELECT * FROM "Anime" ORDER BY "ratingCount" DESC LIMIT 14`
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[MangaRepo])
}

func (q *mangaRepository) UpdateMangaPopularity(ctx context.Context, name string) error {
	query := `UPDATE "Anime" SET "popularity" = "popularity" + 1 WHERE "name" = $1`
	_, err := q.db.Exec(ctx, query, name)
	return err
}
