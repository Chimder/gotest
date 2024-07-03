package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chimas/GoProject/internal/auth"
	"github.com/chimas/GoProject/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type SuccessResponse struct {
	Success string `json:"success"`
}

type User struct {
	Id        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Favorite  pq.StringArray `json:"favorite"`
	CreatedAt time.Time      `json:"createdAt" db:"createdAt"`
}

func NewUserHandler(pgdb *sqlx.DB, rdb *redis.Client) *UserHandler {
	return &UserHandler{pgdb: pgdb, rdb: rdb}
}

type UserHandler struct {
	pgdb *sqlx.DB
	rdb  *redis.Client
}

// @Summary Get a user by email
// @Description Retrieve a user its email
// @Tags User
// @ID get-user-by-email
// @Accept  json
// @Produce  json
// @Param  email path string true "User Email"
// @Success 200 {object} UserSwag
// @Router /user/{email} [get]
func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	op := "handler GetUser"
	email := r.PathValue("email")
	var user User

	err := u.pgdb.Get(&user, `SELECT * FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}

	if err := utils.WriteJSON(w, 200, user); err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}
}

type FavoriteResponse struct {
	IsFavorite bool `json:"isFavorite"`
}
type MangasSwags struct {
	Mangas []MangaSwag `json:"Mangas"`
}

// @Summary User favorite Mangas
// @Description User Favorites
// @Tags User
// @ID get-user-list-manga
// @Accept  json
// @Produce  json
// @Param  email query string true "email"
// @Success 200 {array} MangaSwag
// @Router /user/favorite/list [get]
func (u *UserHandler) UserFavList(w http.ResponseWriter, r *http.Request) {
	op := "handler UserFavList"
	var user User
	email := r.URL.Query().Get("email")

	err := u.pgdb.Get(&user, `SELECT "favorite" FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, 200, []Manga{})
			return
		}
		utils.WriteError(w, 500, op, err)
		return
	}

	if len(user.Favorite) == 0 {
		utils.WriteJSON(w, 200, []Manga{})
		return
	}

	query := `SELECT * FROM "Anime" WHERE "name" = ANY($1)`
	var favoriteMangas []Manga
	err = u.pgdb.Select(&favoriteMangas, query, pq.Array(user.Favorite))
	if err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}

	if err := utils.WriteJSON(w, 200, favoriteMangas); err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}
}

// @Summary User favorite Manga
// @Description User Favorite
// @Tags User
// @ID get-user-favorite-manga
// @Accept  json
// @Produce  json
// @Param  email query string true "email"
// @Param  name query string true "name"
// @Success 200 {object} FavoriteResponse
// @Router /user/favorite/one [get]
func (u *UserHandler) IsUserFavorite(w http.ResponseWriter, r *http.Request) {
	op := "handler IsUserFavorite"
	var user User
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")

	err := u.pgdb.Get(&user, `SELECT * FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, 500, op, err)
			return
		}
		utils.WriteError(w, 500, op, err)
		return
	}

	isAnimeInFavorites := false
	for _, favorite := range user.Favorite {
		if favorite == name {
			isAnimeInFavorites = true
			break
		}
	}
	log.Println("Is Fav:", isAnimeInFavorites)
	w.WriteHeader(http.StatusOK)
	if err := utils.WriteJSON(w, 200, FavoriteResponse{IsFavorite: isAnimeInFavorites}); err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}
}

// @Summary delete user by email
// @Description Delete user
// @Tags User
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param  email query string true "email"
// @Success 200 {object} SuccessResponse
// @Router /user/delete [delete]
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	op := "handler DeleteUser"
	email := r.URL.Query().Get("email")

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		utils.WriteError(w, 401, op, nil)
		return
	}

	if user.Email != email {
		utils.WriteError(w, 403, "Email does not match", nil)
		return
	}

	result, err := u.pgdb.Exec(`DELETE FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}

	if rowsAffected == 0 {
		utils.WriteError(w, 404, "User not found", nil)
		return
	}

	if err := utils.WriteJSON(w, 200, SuccessResponse{Success: "User deleted"}); err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}
}

// @Summary Create or cheack user
// @Description Create
// @Tags User
// @ID create-or-cheack-user
// @Accept  json
// @Produce  json
// @Param  body body string true "Auth Body"
// @Success 200 {object} UserSwag
// @Router /user/create [post]
func (u *UserHandler) CreateOrCheckUser(w http.ResponseWriter, r *http.Request) {
	op := "handler CreateOrCheckUser"
	log.Println("start")
	var newUser User
	if err := utils.ParseJSON(r, &newUser); err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}

	log.Println("FirstUSERRRER", newUser)
	err := u.pgdb.Get(&newUser, `SELECT * FROM "User" WHERE "email" = $1`, newUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			query := `INSERT INTO "User" (id, email, name, image) VALUES (:id, :email, :name, :image)`
			log.Println("INSERT INTO")
			_, err = u.pgdb.NamedExec(query, newUser)
			if err != nil {
				utils.WriteError(w, 500, op, err)
				return
			}
		} else {
			utils.WriteError(w, 500, op, err)
			return
		}
	}
	log.Println("SEcondUSERRRER", newUser)

	log.Println("encrypt")
	encoded, err := auth.Encrypt(newUser)
	if err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}
	log.Println("Cookie sets")
	cookie := &http.Cookie{
		Name:     "manka_google_user",
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	log.Println("seted Cookei", cookie)
	http.SetCookie(w, cookie)

	if err := utils.WriteJSON(w, 200, []byte("register Ok")); err != nil {
		utils.WriteError(w, 500, op, err)
		return
	}
}

// @Summary Get User Session
// @Description Get User Session
// @Tags User
// @ID get-user-session
// @Accept  json
// @Produce  json
// @Success 200 {object} UserSwag
// @Router /user/session [get]
func (u *UserHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	op := "handler GetSession"
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		utils.WriteError(w, 401, op, nil)
		return
	}

	if err := utils.WriteJSON(w, 200, user); err != nil {
		utils.WriteError(w, 401, op, err)
	}
}

// @Summary Toggle Favorite manga
// @Description Toggle manga
// @Tags User
// @ID toggle-favorite-manga
// @Accept  json
// @Produce  json
// @Param  name query string true "manga name"
// @Param  email query string true "email"
// @Success 200 {object} SuccessResponse
// @Router /user/toggle/favorite [post]
func (u *UserHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	op := "handler ToggleFavorite"
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")

	var user User

	err := u.pgdb.Get(&user, `SELECT * FROM "User" WHERE "email" = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, 500, op, err)
			return
		}
		utils.WriteError(w, 500, op, err)
		return
	}

	isAnimeInFavorites := false
	for _, favorite := range user.Favorite {
		if favorite == name {
			isAnimeInFavorites = true
			break
		}
	}

	if !isAnimeInFavorites {

		_, err = u.pgdb.Exec(`UPDATE "Anime" SET "popularity" = popularity + 1 WHERE "name" = $1`, name)
		if err != nil {
			utils.WriteError(w, 500, op, err)
			return
		}

		user.Favorite = append(user.Favorite, name)
		_, err = u.pgdb.NamedExec(`UPDATE "User" SET "favorite" = :favorite WHERE "email" = :email`, user)
		if err != nil {
			utils.WriteError(w, 500, op, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(SuccessResponse{Success: "Manga added"}); err != nil {
			utils.WriteError(w, 500, op, err)
			return
		}
	} else {
		newFavorites := []string{}
		for _, favorite := range user.Favorite {
			if favorite != name {
				newFavorites = append(newFavorites, favorite)
			}
		}
		user.Favorite = newFavorites
		_, err = u.pgdb.NamedExec(`UPDATE "User" SET "favorite" = :favorite WHERE "email" = :email`, user)
		if err != nil {
			utils.WriteError(w, 500, op, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := utils.WriteJSON(w, 200, SuccessResponse{Success: "Manga delete"}); err != nil {
			utils.WriteError(w, 500, op, err)
			return
		}
	}
}
