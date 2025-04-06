package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chimas/GoProject/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MangaRepository interface {
	GetMangaByName(ctx context.Context, name string) (models.MangaRepo, error)
	GetMangaByNames(ctx context.Context, names []string) ([]models.MangaRepo, error)
	ListMangas(ctx context.Context) ([]models.MangaRepo, error)
	ListPopularMangas(ctx context.Context) ([]models.MangaRepo, error)
	UpdateMangaPopularity(ctx context.Context, name string) error
	FilterMangas(ctx context.Context, f models.MangaFilter) ([]models.MangaRepo, error)
}

type mangaRepository struct {
	db *pgxpool.Pool
}

func NewMangaRepository(db *pgxpool.Pool) MangaRepository {
	return &mangaRepository{
		db: db,
	}
}

func (q *mangaRepository) GetMangaByName(ctx context.Context, name string) (models.MangaRepo, error) {
	query := `SELECT * FROM "Anime" WHERE name = $1`
	rows, err := q.db.Query(ctx, query, name)
	if err != nil {
		return models.MangaRepo{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[models.MangaRepo])
}

func (q *mangaRepository) ListMangas(ctx context.Context) ([]models.MangaRepo, error) {
	query := `SELECT * FROM "Anime"`
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MangaRepo])
}

func (q *mangaRepository) GetMangaByNames(ctx context.Context, names []string) ([]models.MangaRepo, error) {
	query := `SELECT * FROM "Anime" WHERE "name" = ANY(@name::text[])`
	rows, err := q.db.Query(ctx, query, pgx.NamedArgs{"name": names})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MangaRepo])
}

func (q *mangaRepository) ListPopularMangas(ctx context.Context) ([]models.MangaRepo, error) {
	query := `SELECT * FROM "Anime" ORDER BY "ratingCount" DESC LIMIT 14`
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MangaRepo])
}

func (q *mangaRepository) FilterMangas(ctx context.Context, f models.MangaFilter) ([]models.MangaRepo, error) {
	args := []interface{}{}
	where := []string{}
	i := 1

	query := strings.Builder{}
	query.WriteString(`SELECT * FROM "Anime"`)
	if f.Name != "" {
		where = append(where, fmt.Sprintf(`"name" ILIKE $%d`, i))
		args = append(args, "%"+f.Name+"%")
		i++
	}
	if f.Status != "" {
		where = append(where, fmt.Sprintf(`"status" = $%d`, i))
		args = append(args, f.Status)
		i++
	}
	if f.Country != "" {
		where = append(where, fmt.Sprintf(`"country" = $%d`, i))
		args = append(args, f.Country)
		i++
	}
	for _, genre := range f.Genres {
		where = append(where, fmt.Sprintf(`"genres" @> ARRAY[$%d]`, i))
		args = append(args, genre)
		i++
	}

	if len(where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(where, " AND "))

	}
	if f.OrderField != "" && f.OrderSort != "" {
		query.WriteString(fmt.Sprintf(` ORDER BY "%s" %s`, f.OrderField, f.OrderSort))
	}
	if f.Page > 0 && f.PerPage > 0 {
		query.WriteString(fmt.Sprintf(` LIMIT %d OFFSET %d`, f.PerPage, (f.Page-1)*f.PerPage))
	}

	rows, err := q.db.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("repo.FilterMangas query: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.MangaRepo])
}

func (q *mangaRepository) UpdateMangaPopularity(ctx context.Context, name string) error {
	query := `UPDATE "Anime" SET "popularity" = "popularity" + 1 WHERE "name" = $1`
	_, err := q.db.Exec(ctx, query, name)
	return err
}
