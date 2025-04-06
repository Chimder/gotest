package models

import (
	"time"
)

type ChapterRepo struct {
	Chapter   int32     `db:"chapter"`
	Img       []string  `db:"img"`
	Name      string    `db:"name"`
	AnimeName string    `db:"animeName"`
	CreatedAt time.Time `db:"createdAt"`
}

type ChapterResp struct {
	Chapter   int       `json:"chapter"`
	Img       []string  `json:"genres"`
	Name      string    `json:"name"`
	MangaName string    `json:"mangaName"`
	CreatedAt time.Time `json:"createdAt"`
}

func ChapterRespFromDB(c ChapterRepo) *ChapterResp {
	return &ChapterResp{
		Chapter:   int(c.Chapter),
		Img:       c.Img,
		Name:      c.Name,
		MangaName: c.AnimeName,
		CreatedAt: c.CreatedAt,
	}
}

func ChaptersRespFromDB(ch []ChapterRepo) []ChapterResp {
	repo := make([]ChapterResp, len(ch))
	for i, r := range ch {
		repo[i] = *ChapterRespFromDB(r)
	}
	return repo
}
