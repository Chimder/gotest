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

func (s *ChapterService) GetChapterByAnimeNameAndNumber(ctx context.Context, name string, chapter int) (repository.ChapterRepo, error) {
	return s.repo.Chapter.GetChapterByAnimeNameAndNumber(ctx, name, chapter)
}

func (s *ChapterService) ListChaptersByAnimeName(ctx context.Context, animename string) ([]repository.ChapterRepo, error) {
	return s.repo.Chapter.ListChaptersByAnimeName(ctx, animename)
}
