package service

import (
	"context"

	"github.com/chimas/GoProject/internal/repository"
)

type MangaService struct {
	repo *repository.Repository
}

func NewMangaService(repo *repository.Repository) *MangaService {
	return &MangaService{repo: repo}
}

type MangaWithChapters struct {
	Manga    repository.MangaRepo
	Chapters []repository.ChapterRepo `json:"chapters"`
}

func (s *MangaService) GetMangaByName(ctx context.Context, name string) (*MangaWithChapters, error) {
	manga, err := s.repo.Manga.GetMangaByName(ctx, name)
	if err != nil {
		return nil, err
	}

	chapters, err := s.repo.Chapter.ListChaptersByAnimeName(ctx, name)
	if err != nil {
		return nil, err

	}

	mangaWithChapter := &MangaWithChapters{
		Manga:    manga,
		Chapters: chapters,
	}

	return mangaWithChapter, err
}

// if err := utils.WriteJSON(w, 200, &mangaWithChapter); err != nil {

// 	return repository.MangaRepo{}, err
// }

// return s.repo.Manga.GetMangaByName(ctx, name)
// }

func (s *MangaService) ListMangas(ctx context.Context) ([]repository.MangaRepo, error) {
	return s.repo.Manga.ListMangas(ctx)
}

func (s *MangaService) ListPopularMangas(ctx context.Context) ([]repository.MangaRepo, error) {
	return s.repo.Manga.ListPopularMangas(ctx)
}

// func (s *MangaService) UpdateMangaPopularity(ctx context.Context, name string) error {
// 	return s.repo.Manga.UpdateMangaPopularity(ctx, name)
// }
