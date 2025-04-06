package service

import (
	"context"

	"github.com/chimas/GoProject/internal/models"
	"github.com/chimas/GoProject/internal/repository"
)

type MangaService struct {
	repo *repository.Repository
}

func NewMangaService(repo *repository.Repository) *MangaService {
	return &MangaService{repo: repo}
}

func (s *MangaService) GetMangaByName(ctx context.Context, name string) (*models.MangaWithChaptersResp, error) {
	manga, err := s.repo.Manga.GetMangaByName(ctx, name)
	if err != nil {
		return nil, err
	}
	chapters, err := s.repo.Chapter.ListChaptersByMangaName(ctx, name)
	if err != nil {
		return nil, err

	}

	resp := models.MangaWithChaptersRespFromDB(manga, chapters)
	return &resp, err
}

func (s *MangaService) ListMangas(ctx context.Context) ([]models.MangaRepo, error) {
	return s.repo.Manga.ListMangas(ctx)
}

func (s *MangaService) ListPopularMangas(ctx context.Context) ([]models.MangaRepo, error) {
	return s.repo.Manga.ListPopularMangas(ctx)
}

func (s *MangaService) FilterMangas(ctx context.Context, filter models.MangaFilter) ([]models.MangaRepo, error) {
	return s.repo.Manga.FilterMangas(ctx, filter)
}

// func (s *MangaService) FilterManga(ctx context.Context, name string) ([]repository.MangaRepo, error) {

// }
