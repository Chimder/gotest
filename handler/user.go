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

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	result, err := u.db.Exec(`DELETE FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "User deleted"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (u *UserHandler) CreateUserIfNotExists(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = u.db.Get(&newUser, `SELECT * FROM "User" WHERE "email" = $1`, newUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			query := `INSERT INTO "User" (id, email, name, image ) VALUES (:id, :email, :name, :image)`
			_, err = u.db.NamedExec(query, newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u *UserHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	var user User
	name := r.PathValue("name")
	email := r.PathValue("email")
	err := u.db.Get(&user, `SELECT * FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	isAnimeInFavorites := false
	for _, favorite := range user.Favorite {
		if favorite == name {
			isAnimeInFavorites = true
			break
		}
	}

	log.Println("sisidisis", isAnimeInFavorites)

	if !isAnimeInFavorites {

		_, err = u.db.Exec(`UPDATE "Anime" SET "popularity" = popularity + 1 WHERE "name" = $1`, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		user.Favorite = append(user.Favorite, name)
		_, err = u.db.NamedExec(`UPDATE "User" SET "favorite" = :favorite WHERE "email" = :email`, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		newFavorites := []string{}
		for _, favorite := range user.Favorite {
			if favorite != name {
				newFavorites = append(newFavorites, favorite)
			}
		}
		user.Favorite = newFavorites
		_, err = u.db.NamedExec(`UPDATE "User" SET "favorite" = :favorite WHERE "email" = :email`, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
