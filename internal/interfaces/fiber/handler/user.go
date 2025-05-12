package handler

import (
	"fmt"
	"time"

	"github.com/chimas/GoProject/internal/auth"
	"github.com/chimas/GoProject/internal/models"
	"github.com/chimas/GoProject/internal/service"
	"github.com/gofiber/fiber/v2"
)

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
// @Success 200 {object} models.UserResp
// @Router /user [get]
func (u *UserHandler) GetUser(c *fiber.Ctx) error {
	op := "handler GetUser"
	email := c.Query("email")

	user, err := u.serv.GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "GUBE",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

type FavoriteResponse struct {
	IsFavorite bool `json:"isFavorite"`
}

// @Summary User favorite Mangas
// @Description User Favorites
// @Tags User
// @ID get-user-list-manga
// @Accept  json
// @Produce  json
// @Param  email query string true "email"
// @Success 200 {array} models.MangaResp
// @Router /user/favorite/list [get]
func (u *UserHandler) UserFavList(c *fiber.Ctx) error {
	op := "handler UserFavList"
	email := c.Query("email")

	favorites, err := u.serv.GetUserFavorites(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "GUFBE",
		})
	}
	return c.Status(fiber.StatusOK).JSON(favorites)
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
func (u *UserHandler) IsUserFavorite(c *fiber.Ctx) error {
	op := "handler IsUserFavorite"
	name := c.Query("name")
	email := c.Query("email")

	isMangaIsFavorite, err := u.serv.IsUserFavorite(c.Context(), email, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "NIL",
		})
	}
	return c.Status(fiber.StatusOK).JSON(FavoriteResponse{IsFavorite: isMangaIsFavorite})
}

type SuccessResponse struct {
	Success string `json:"success"`
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
func (u *UserHandler) DeleteUser(c *fiber.Ctx) error {
	op := "handler DeleteUser"
	email := c.Query("email")

	user, ok := auth.GetUserFromContext(c.Context())
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": op + "GUFC",
		})
	}

	if user.Email != email {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Email does not match",
		})
	}

	err := u.serv.DeleteUser(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "DUBE",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "manka_google_user",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{Success: "User deleted"})
}

// @Summary Create or check user
// @Description Create
// @Tags User
// @ID create-or-check-user
// @Accept  json
// @Produce  json
// @Param  body body string true "Auth Body"
// @Success 200 {object} models.UserResp
// @Router /user/create [post]
func (u *UserHandler) CreateOrCheckUser(c *fiber.Ctx) error {
	op := "handler CreateOrCheckUser"
	var newUser models.UserRepo

	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "PJ",
		})
	}

	user, err := u.serv.InsertUser(c.Context(), &newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "Err insert user",
		})
	}

	encoded, err := auth.Encrypt(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "ENC",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "manka_google_user",
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary DeleteUserCookie
// @Description delete user cookie
// @Tags User
// @ID delete-user-cookie
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Empty
// @Router /user/cookie/delete [get]
func (u *UserHandler) DeleteCookie(c *fiber.Ctx) error {
	op := "handle/DeleteCookie"
	cookieName := "manka_google_user"

	if cookieValue := c.Cookies(cookieName); cookieValue == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": op + "CK",
		})
	}

	c.ClearCookie(cookieName)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("%s: cookie %s deleted", op, cookieName),
	})
}

// @Summary Get User Session
// @Description Get User Session
// @Tags User
// @ID get-user-session
// @Accept  json
// @Produce  json
// @Success 200 {object} models.UserResp
// @Router /user/session [get]
func (u *UserHandler) GetSession(c *fiber.Ctx) error {
	op := "handler GetSession"
	user, ok := auth.GetUserFromContext(c.Context())
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": op + "GUFC",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
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
func (u *UserHandler) ToggleFavorite(c *fiber.Ctx) error {
	op := "handler ToggleFavorite"
	name := c.Query("name")
	email := c.Query("email")

	err := u.serv.ToggleFavorite(c.Context(), email, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "UUF2",
		})
	}
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{Success: "Manga toggled"})
}
