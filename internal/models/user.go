package models

import (
	"time"

	"github.com/jackc/pgx/pgtype"
)

type UserRepo struct {
	ID        string           `db:"id"`
	Name      string           `db:"name"`
	Email     string           `db:"email"`
	Image     string           `db:"image"`
	Favorite  []string         `db:"favorite"`
	CreatedAt pgtype.Timestamp `db:"createdAt"`
}

type UserResp struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Favorite  []string  `json:"favorite"`
	CreatedAt time.Time `json:"createdAt"`
}

// type User struct {
// 	Id        string         `json:"id"`
// 	Email     string         `json:"email"`
// 	Name      string         `json:"name"`
// 	Image     string         `json:"image"`
// 	Favorite  pq.StringArray `json:"favorite"`
// 	CreatedAt time.Time      `json:"createdAt" db:"createdAt"`
// }
