package handler

import "time"

type MangaSwag struct {
	Name          string        `json:"name"`
	Img           string        `json:"img"`
	ImgHeader     string        `json:"imgHeader" db:"imgHeader"`
	Describe      string        `json:"describe"`
	Genres        []string      `json:"genres" db:"genres"`
	Author        string        `json:"author"`
	Country       string        `json:"country"`
	Published     int           `json:"published"`
	AverageRating float64       `json:"averageRating" db:"averageRating"`
	RatingCount   int           `json:"ratingCount" db:"ratingCount"`
	Status        string        `json:"status"`
	Popularity    int           `json:"popularity"`
	Id            int           `json:"id"`
	Chapters      []ChapterSwag `json:"chapters"`
}

type ChapterSwag struct {
	Chapter   int       `json:"chapter"`
	Img       []string  `json:"genres" db:"img"`
	Name      string    `json:"name"`
	AnimeName string    `json:"animeName" db:"animeName"`
	CreatedAt time.Time `json:"created" db:"createdAt"`
}
