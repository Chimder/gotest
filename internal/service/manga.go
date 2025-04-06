package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/chimas/GoProject/internal/models"
	"github.com/chimas/GoProject/internal/repository"
	"github.com/redis/go-redis/v9"
)

type MangaService struct {
	repo *repository.Repository
	rdb  *redis.Client
}

func NewMangaService(repo *repository.Repository, rdb *redis.Client) *MangaService {
	return &MangaService{repo: repo, rdb: rdb}
}

func (s *MangaService) GetMangaByName(ctx context.Context, name string) ([]byte, error) {
	key := "MangaByName" + name

	if manga, err := s.rdb.Get(ctx, key).Bytes(); err == nil {
		return manga, nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var manga models.MangaRepo
	var chapters []models.ChapterRepo
	var mangaErr, chaptersErr error

	start := time.Now()

	go func() {
		defer wg.Done()
		manga, mangaErr = s.repo.Manga.GetMangaByName(ctx, name)
	}()

	go func() {
		defer wg.Done()
		chapters, chaptersErr = s.repo.Chapter.ListChaptersByMangaName(ctx, name)
	}()

	wg.Wait()
	parallelTime := time.Since(start)
	slog.Info("Execution times", "parallel", parallelTime)

	if mangaErr != nil {
		return nil, fmt.Errorf("failed to get manga: %w", mangaErr)
	}
	if chaptersErr != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", chaptersErr)
	}

	resp := models.MangaWithChaptersRespFromDB(manga, chapters)
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err

	}

	go func() {
		if err := s.rdb.Set(context.Background(), key, data, 1*time.Minute).Err(); err != nil {
			slog.Error("failed to cache manga", "key", key, "error", err)
		}
	}()

	return data, err
}
func (s *MangaService) ListMangas(ctx context.Context) ([]byte, error) {
	key := "ListMangas"

	if data, err := s.rdb.Get(ctx, key).Bytes(); err == nil {
		return data, nil
	}
	mangas, err := s.repo.Manga.ListMangas(ctx)
	if err != nil {
		return nil, err
	}

	resp := models.MangasRespFromDB(mangas)
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.rdb.Set(context.Background(), key, data, 1*time.Minute).Err(); err != nil {
			slog.Error("failed to cache manga", "key", key, "error", err)
		}
	}()

	return data, err
}

func (s *MangaService) ListPopularMangas(ctx context.Context) ([]byte, error) {
	key := "ListPopularMangas"

	if data, err := s.rdb.Get(ctx, key).Bytes(); err == nil {
		return data, nil
	}

	mangas, err := s.repo.Manga.ListPopularMangas(ctx)
	if err != nil {
		return nil, err
	}

	resp := models.MangasRespFromDB(mangas)
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.rdb.Set(context.Background(), key, data, 1*time.Minute).Err(); err != nil {
			slog.Error("failed to cache manga", "key", key, "error", err)
		}
	}()

	return data, err
}

func (s *MangaService) FilterMangas(ctx context.Context, filter models.MangaFilter) ([]models.MangaRepo, error) {
	return s.repo.Manga.FilterMangas(ctx, filter)
}
