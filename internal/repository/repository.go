package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	User    UserRepository
	Manga   MangaRepository
	Chapter ChapterRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Manga:   NewMangaRepository(db),
		Chapter: NewChapterRepository(db),
	}
}
