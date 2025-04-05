package handler

import "time"

type ChapterSwag struct {
	Chapter   int       `json:"chapter"`
	Img       []string  `json:"genres" db:"img"`
	Name      string    `json:"name"`
	AnimeName string    `json:"animeName" db:"animeName"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}
