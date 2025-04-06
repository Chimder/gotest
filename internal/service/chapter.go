package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/chimas/GoProject/internal/models"
	"github.com/chimas/GoProject/internal/repository"
	"github.com/redis/go-redis/v9"
)

type ChapterService struct {
	repo *repository.Repository
	rdb  *redis.Client
}

func NewChapterService(repo *repository.Repository, rdb *redis.Client) *ChapterService {
	return &ChapterService{repo: repo, rdb: rdb}
}

func (s *ChapterService) GetChapterByMangaNameAndNumber(ctx context.Context, name string, chapter int) ([]byte, error) {
	key := "GetChapterByManga" + name + strconv.Itoa(chapter)

	if data, err := s.rdb.Get(ctx, key).Bytes(); err == nil {
		return data, nil
	}

	chapters, err := s.repo.Chapter.GetChapterByMangaNameAndNumber(ctx, name, chapter)
	if err != nil {
		return nil, err
	}

	resp := models.ChapterRespFromDB(chapters)
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

func (s *ChapterService) ListChaptersByMangaName(ctx context.Context, name string) ([]models.ChapterRepo, error) {
	return s.repo.Chapter.ListChaptersByMangaName(ctx, name)
}
