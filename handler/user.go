package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type User struct {
	Id        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Favorite  pq.StringArray `json:"favorite"`
	CreatedAt time.Time      `json:"createdAt" db:"createdAt"`
}

func NewUserHandler(db *sqlx.DB, rdb *redis.Client) *UserHandler {
	return &UserHandler{db: db, rdb: rdb}
}

type UserHandler struct {
	db  *sqlx.DB
	rdb *redis.Client
}

// func (u *UserHandler) GetUser(email string) (*User, error) {
func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	email := r.PathValue("email")
	log.Println(email)
	var user User

	err := u.db.Get(&user, `SELECT * FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u *UserHandler) CreateUserIfNotExists(w http.ResponseWriter, r *http.Request) error {
	email := r.PathValue("email")

	var user User
	err := u.db.Get(&user, `SELECT * FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Пользователь не найден, создаем нового
			query := `INSERT INTO User (id, email, name, image, favorite, createdAt) VALUES (:id, :email, :name, :image, :favorite, :createdAt)`
			_, err = u.db.NamedExec(query, newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// Пользователь уже существует, пропускаем
	return nil
}
