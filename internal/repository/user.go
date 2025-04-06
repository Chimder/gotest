package repository

import (
	"context"

	"github.com/chimas/GoProject/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (models.UserRepo, error)
	GetUserFavoritesByEmail(ctx context.Context, email string) ([]string, error)
	UpdateUserFavorites(ctx context.Context, favorite []string, email string) error
	DeleteUserByEmail(ctx context.Context, email string) error
	InsertUser(ctx context.Context, arg InsertUserParams) (models.UserRepo, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (q *userRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, `DELETE FROM "User" WHERE "email" = $1`, email)
	return err
}

func (q *userRepository) GetUserByEmail(ctx context.Context, email string) (models.UserRepo, error) {
	query := `SELECT * FROM "User" WHERE "email" = $1;`
	rows, err := q.db.Query(ctx, query, email)
	if err != nil {
		return models.UserRepo{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserRepo])
}

func (q *userRepository) GetUserFavoritesByEmail(ctx context.Context, email string) ([]string, error) {
	query := `SELECT "favorite" FROM "User" WHERE "email" = $1`
	var favorite []string
	err := q.db.QueryRow(ctx, query, email).Scan(&favorite)
	return favorite, err
}

type InsertUserParams struct {
	ID    string `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
	Name  string `db:"name" json:"name"`
	Image string `db:"image" json:"image"`
}

func (q *userRepository) InsertUser(ctx context.Context, arg InsertUserParams) (models.UserRepo, error) {
	query := `
INSERT INTO "User" (id, email, name, image)
VALUES (@id, @email, @name, @image)
RETURNING *
	`
	rows, err := q.db.Query(ctx, query, pgx.NamedArgs{
		"id":    arg.ID,
		"email": arg.Email,
		"name":  arg.Name,
		"image": arg.Image})
	if err != nil {
		return models.UserRepo{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.UserRepo])
}

func (q *userRepository) UpdateUserFavorites(ctx context.Context, favorite []string, email string) error {
	query := `UPDATE "User" SET "favorite" = $1 WHERE "email" = $2`
	_, err := q.db.Exec(ctx, query, favorite, email)
	return err
}
