package handler

import (
	"net/http"
	"strconv"

	"github.com/chimas/GoProject/internal/service"
	"github.com/chimas/GoProject/utils"
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
// @Success 200 {object} ChapterSwag
// @Router /manga/chapter [get]
func (m *ChapterHandler) Chapter(w http.ResponseWriter, r *http.Request) {
	op := "handler Chapter"
	name := r.URL.Query().Get("name")
	chapStr := r.URL.Query().Get("chapter")

	chap, err := strconv.Atoi(chapStr)
	if err != nil {
		utils.WriteError(w, 400, op+"ATOI", err)
		return
	}
	chapter, err := m.serv.GetChapterByMangaNameAndNumber(r.Context(), name, chap)
	if err != nil {
		utils.WriteError(w, 500, op+"GCBANAN", err)
		return
	}

	if err := utils.WriteJSON(w, 200, &chapter); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
		return
	}
}
