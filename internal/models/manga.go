package models

import "github.com/jackc/pgx/pgtype"

type MangaRepo struct {
	ID            int32       `db:"id"`
	Name          string      `db:"name"`
	Img           string      `db:"img"`
	ImgHeader     string      `db:"imgHeader"`
	Describe      string      `db:"describe"`
	Genres        []string    `db:"genres"`
	Author        pgtype.Text `db:"author"`
	Country       string      `db:"country"`
	Published     int         `db:"published"`
	AverageRating float32     `db:"averageRating"`
	RatingCount   int         `db:"ratingCount"`
	Status        string      `db:"status"`
	Popularity    int         `db:"popularity"`
}

type MangaWithChaptersResp struct {
	Manga    *MangaResp    `json:"manga"`
	Chapters []ChapterResp `json:"chapters"`
}

func MangaWithChaptersRespFromDB(m MangaRepo, c []ChapterRepo) MangaWithChaptersResp {
	return MangaWithChaptersResp{
		Manga:    MangaRespFromDB(m),
		Chapters: ChaptersRespFromDB(c),
	}
}

func MangaRespFromDB(m MangaRepo) *MangaResp {
	return &MangaResp{
		Id:            int(m.ID),
		Name:          m.Name,
		Img:           m.Img,
		ImgHeader:     m.ImgHeader,
		Describe:      m.Describe,
		Genres:        m.Genres,
		Author:        m.Author.String,
		Country:       m.Country,
		Published:     m.Published,
		AverageRating: float64(m.AverageRating),
		RatingCount:   m.RatingCount,
		Status:        m.Status,
		Popularity:    m.Popularity,
	}
}

func MangasRespFromDB(m []MangaRepo) []MangaResp {
	repo := make([]MangaResp, len(m))
	for i, r := range m {
		repo[i] = *MangaRespFromDB(r)
	}
	return repo
}

type MangaFilter struct {
	Name       string
	Genres     []string
	Status     string
	Country    string
	OrderField string
	OrderSort  string
	Page       int
	PerPage    int
}

type MangaResp struct {
	Name          string   `json:"name"`
	Img           string   `json:"img"`
	ImgHeader     string   `json:"imgHeader"`
	Describe      string   `json:"describe"`
	Genres        []string `json:"genres"`
	Author        string   `json:"author"`
	Country       string   `json:"country"`
	Published     int      `json:"published"`
	AverageRating float64  `json:"averageRating"`
	RatingCount   int      `json:"ratingCount"`
	Status        string   `json:"status"`
	Popularity    int      `json:"popularity"`
	Id            int      `json:"id"`
	// Chapters      []ChapterResp `json:"chapters"`
}

type Empty struct {
}

type MangasResp struct {
	Mangas []MangaResp `json:"Mangas"`
}
