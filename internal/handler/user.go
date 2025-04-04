package handler

import (
	"net/http"
	"time"

	"github.com/chimas/GoProject/internal/auth"
	"github.com/chimas/GoProject/internal/repository"
	"github.com/chimas/GoProject/internal/service"
	"github.com/chimas/GoProject/utils"
	"github.com/lib/pq"
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
type UserHandler struct {
	serv *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{serv: s}
}

// @Summary Get a user by email
// @Description Retrieve a user its email
// @Tags User
// @ID get-user-by-email
// @Accept  json
// @Produce  json
// @Param  email query string true "User Email"
// @Success 200 {object} UserSwag
// @Router /user [get]
func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	op := "handler GetUser"
	email := r.URL.Query().Get("email")

	user, err := u.serv.GetUserByEmail(r.Context(), email)
	if err != nil {
		utils.WriteError(w, 500, op+"GUBE", err)
		return
	}

	if err := utils.WriteJSON(w, 200, &user); err != nil {
		utils.WriteError(w, 500, op+"WJ2", err)
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
	email := r.URL.Query().Get("email")

	favorites, err := u.serv.GetUserFavorites(r.Context(), email)
	if err != nil {
		utils.WriteError(w, 500, op+"GUFBE", err)
		return
	}

	if err := utils.WriteJSON(w, 200, &favorites); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
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
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")

	isMangaIsFavorite, err := u.serv.IsUserFavorite(r.Context(), email, name)
	if err != nil {
		utils.WriteError(w, 500, op+"NIL", err)
		return
	}

	if err := utils.WriteJSON(w, 200, FavoriteResponse{IsFavorite: isMangaIsFavorite}); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
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
		utils.WriteError(w, 401, op+"GUFC", nil)
		return
	}

	if user.Email != email {
		utils.WriteError(w, 403, "Email does not match", nil)
		return
	}

	err := u.serv.DeleteUser(r.Context(), email)
	if err != nil {
		utils.WriteError(w, 500, op+"DUBE", err)
		return
	}

	cookie := &http.Cookie{
		Name:     "manka_google_user",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	if err := utils.WriteJSON(w, 200, SuccessResponse{Success: "User deleted"}); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
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
	var newUser repository.UserRepo

	if err := utils.ParseJSON(r, &newUser); err != nil {
		utils.WriteError(w, 500, op+"PJ", err)
		return
	}

	user, err := u.serv.InsertUser(r.Context(), &newUser)
	if err != nil {
		utils.WriteError(w, 500, op+"Err insert user", err)
		return
	}

	encoded, err := auth.Encrypt(user)
	if err != nil {
		utils.WriteError(w, 500, op+"ENC", err)
		return
	}

	cookie := &http.Cookie{
		Name:     "manka_google_user",
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)

	if err := utils.WriteJSON(w, 200, &user); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
		return
	}
}

// @Summary DeleteUserCookie
// @Description delete user cookie
// @Tags User
// @ID delete-user-cookie
// @Accept  json
// @Produce  json
// @Success 200 {array} Empty
// @Router /user/cookie/delete [get]
func (u *UserHandler) DeleteCookie(w http.ResponseWriter, r *http.Request) {
	op := "handle/DeleteCookie"
	cookieName := "manka_google_user"

	_, err := r.Cookie(cookieName)
	if err != nil {
		utils.WriteError(w, 404, op+"CK", err)
		return
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("cookie deleted"))
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
		utils.WriteError(w, 401, op+"GUFC", nil)
		return
	}

	if err := utils.WriteJSON(w, 200, &user); err != nil {
		utils.WriteError(w, 401, op+"WJ", err)
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

	err := u.serv.ToggleFavorite(r.Context(), email, name)
	if err != nil {
		utils.WriteError(w, 500, op+"UUF2", err)
		return
	}

	if err := utils.WriteJSON(w, 200, SuccessResponse{Success: "Manga toggled"}); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
		return
	}
}
