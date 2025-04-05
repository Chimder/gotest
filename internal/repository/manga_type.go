package repository

import "github.com/jackc/pgx/pgtype"

type MangaRepo struct {
	ID            int32         `db:"id" json:"id"`
	Name          string        `db:"name" json:"name"`
	Img           string        `db:"img" json:"img"`
	ImgHeader     string        `db:"imgHeader" json:"imgHeader"`
	Describe      string        `db:"describe" json:"describe"`
	Genres        []string      `db:"genres" json:"genres"`
	Author        pgtype.Text   `db:"author" json:"author"`
	Country       string        `db:"country" json:"country"`
	Published     int32         `db:"published" json:"published"`
	AverageRating pgtype.Float8 `db:"averageRating" json:"averageRating"`
	RatingCount   pgtype.Int4   `db:"ratingCount" json:"ratingCount"`
	Status        string        `db:"status" json:"status"`
	Popularity    pgtype.Int4   `db:"popularity" json:"popularity"`
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