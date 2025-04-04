package service

import (
	"context"

	"github.com/chimas/GoProject/internal/repository"
)

type ChapterService struct {
	repo *repository.Repository
}

func NewChapterService(repo *repository.Repository) *ChapterService {
	return &ChapterService{repo: repo}
}

func (s *ChapterService) GetChapterByMangaNameAndNumber(ctx context.Context, name string, chapter int) (repository.ChapterRepo, error) {
	return s.repo.Chapter.GetChapterByMangaNameAndNumber(ctx, name, chapter)
}

func (s *ChapterService) ListChaptersByMangaName(ctx context.Context, name string) ([]repository.ChapterRepo, error) {
	return s.repo.Chapter.ListChaptersByMangaName(ctx, name)
}
