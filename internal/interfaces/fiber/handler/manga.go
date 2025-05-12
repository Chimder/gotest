package handler

import (
	"github.com/chimas/GoProject/internal/service"
	"github.com/gofiber/fiber/v2"
)

type MangaHandler struct {
	serv *service.MangaService
}

func NewMangaHandler(s *service.MangaService) *MangaHandler {
	return &MangaHandler{serv: s}
}

// @Summary Get all mangas
// @Description Retrieve a list of all mangas
// @Tags Manga
// @ID get-all-mangas
// @Accept  json
// @Produce  json
// @Success 200 {array} models.MangaResp
// @Router /manga/many [get]
func (m *MangaHandler) Mangas(c *fiber.Ctx) error {
	op := "handler Mangas"

	mangas, err := m.serv.ListMangas(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "LM",
		})
	}
	return c.Status(fiber.StatusOK).Type("json").Send(mangas)
}

// @Summary Get a manga by name
// @Description Retrieve a manga by its name
// @Tags Manga
// @ID get-manga-by-name
// @Accept  json
// @Produce  json
// @Param  name query string true "Name of the Manga"
// @Success 200 {object} models.MangaResp
// @Router /manga [get]
func (m *MangaHandler) Manga(c *fiber.Ctx) error {
	op := "handler Manga"
	name := c.Query("name")

	manga, err := m.serv.GetMangaByName(c.Context(), name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "GMBN",
		})
	}
	return c.Status(fiber.StatusOK).Type("json").Send(manga)
}

// @Summary Get popular mangas
// @Description Retrieve a list of popular mangas
// @Tags Manga
// @ID get-popular-manga
// @Accept  json
// @Produce  json
// @Success 200 {array} models.MangaResp
// @Router /manga/popular [get]
func (m *MangaHandler) Popular(c *fiber.Ctx) error {
	mangas, err := m.serv.ListPopularMangas(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Get Popular Mangas",
		})
	}
	// return c.Status(fiber.StatusOK).JSON(manga)
	return c.Status(fiber.StatusOK).Type("json").Send(mangas)
}

type FilterParams struct {
	Name       string   `schema:"name"`
	Genres     []string `schema:"genres"`
	Status     string   `schema:"status"`
	Country    string   `schema:"country"`
	OrderField string   `schema:"orderField"`
	OrderSort  string   `schema:"orderSort"`
	Page       int      `schema:"page"`
	PerPage    int      `schema:"perPage"`
}

// @Summary Get a chapter
// @Description Find Manga Chapter
// @Tags Manga
// @ID Filter-anime
// @Accept  json
// @Produce  json
// @Param  name query string false "Name of the Manga"
// @Param  genres[] query []string false "Chapter of the Manga"
// @Param  status query string false "Name of the Manga"
// @Param  country query string false "Chapter of the Manga"
// @Param  orderField query string false "field of the Manga"
// @Param  orderSort query string false "sort of the Manga"
// @Param  page query int false "page not 0"
// @Param  perPage query int false "perPage"
// @Success 200 {array} models.MangaResp
// @Router /manga/filter [get]
// func (m *MangaHandler) Filter(c *fiber.Ctx) error {
// 	page, _ := strconv.Atoi(c.Query("page"))
// 	perPage, _ := strconv.Atoi(c.Query("perPage"))

// 	filter := models.MangaFilter{
// 		Name:       c.Query("name"),
// 		Genres:     c.Context().QueryArgs().PeekMulti("genres[]"),
// 		Status:     c.Query("status"),
// 		Country:    c.Query("country"),
// 		OrderField: c.Query("orderField"),
// 		OrderSort:  c.Query("orderSort"),
// 		Page:       page,
// 		PerPage:    perPage,
// 	}

// 	mangas, err := m.serv.FilterMangas(c.Context(), &filter)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": fmt.Sprintf("filter mangas: %v", err),
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(mangas)
// }
