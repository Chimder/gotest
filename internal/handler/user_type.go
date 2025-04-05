package handler

import (
	"time"

	"github.com/lib/pq"
)

type UserSwag struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Favorite  []string  `json:"favorite"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}

type User struct {
	Id        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Favorite  pq.StringArray `json:"favorite"`
	CreatedAt time.Time      `json:"createdAt" db:"createdAt"`
}
