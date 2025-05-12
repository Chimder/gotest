package handler

import (
	"strconv"

	"github.com/chimas/GoProject/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ChapterHandler struct {
	serv *service.ChapterService
}

func NewChapterHandler(s *service.ChapterService) *ChapterHandler {
	return &ChapterHandler{serv: s}
}

// @Summary Get a chapter
// @Description Find Manga Chapter
// @Tags Chapter
// @ID get-chapter
// @Accept  json
// @Produce  json
// @Param  name query string true "Name of the Manga"
// @Param  chapter query int true "Chapter of the Manga"
// @Success 200 {object} models.ChapterResp
// @Router /manga/chapter [get]
func (m *ChapterHandler) Chapter(c *fiber.Ctx) error {
	op := "handler Chapter"

	name := c.Query("name")
	chapStr := c.Query("chapter")

	chap, err := strconv.Atoi(chapStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": op + "ATOI",
		})
	}

	chapter, err := m.serv.GetChapterByMangaNameAndNumber(c.Context(), name, chap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": op + "GCBANAN",
		})
	}

	return c.Status(fiber.StatusOK).JSON(chapter)
}
