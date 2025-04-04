package repository

import (
	"context"

	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (UserRepo, error)
	GetUserFavoritesByEmail(ctx context.Context, email string) ([]string, error)
	UpdateUserFavorites(ctx context.Context, favorite []string, email string) error
	DeleteUserByEmail(ctx context.Context, email string) error
	InsertUser(ctx context.Context, arg InsertUserParams) (UserRepo, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

type UserRepo struct {
	ID        string           `db:"id" json:"id"`
	Name      string           `db:"name" json:"name"`
	Email     string           `db:"email" json:"email"`
	Image     string           `db:"image" json:"image"`
	Favorite  []string         `db:"favorite" json:"favorite"`
	CreatedAt pgtype.Timestamp `db:"createdAt" json:"createdAt"`
}

func (q *userRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, `DELETE FROM "User" WHERE "email" = $1`, email)
	return err
}

func (q *userRepository) GetUserByEmail(ctx context.Context, email string) (UserRepo, error) {
	query := `SELECT * FROM "User" WHERE "email" = $1;`
	rows, err := q.db.Query(ctx, query, email)
	if err != nil {
		return UserRepo{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[UserRepo])
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

func (q *userRepository) InsertUser(ctx context.Context, arg InsertUserParams) (UserRepo, error) {
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
		return UserRepo{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[UserRepo])
}

func (q *userRepository) UpdateUserFavorites(ctx context.Context, favorite []string, email string) error {
	query := `UPDATE "User" SET "favorite" = $1 WHERE "email" = $2`
	_, err := q.db.Exec(ctx, query, favorite, email)
	return err
}
