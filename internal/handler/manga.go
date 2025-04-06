package handler

import (
	"net/http"
	"strconv"

	"github.com/chimas/GoProject/internal/models"
	"github.com/chimas/GoProject/internal/service"
	"github.com/chimas/GoProject/utils"
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
func (m *MangaHandler) Mangas(w http.ResponseWriter, r *http.Request) {
	op := "handler Mangas"

	mangas, err := m.serv.ListMangas(r.Context())
	if err != nil {
		utils.WriteError(w, 500, op+"LM")
		return
	}
	resp := models.MangasRespFromDB(mangas)
	utils.WriteJSON(w, 200, &resp)
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
func (m *MangaHandler) Manga(w http.ResponseWriter, r *http.Request) {
	op := "handler Manga"
	name := r.URL.Query().Get("name")

	manga, err := m.serv.GetMangaByName(r.Context(), name)
	if err != nil {
		utils.WriteError(w, 500, op+"GMBN")
		return
	}

	utils.WriteJSON(w, 200, &manga)
}

// @Summary Get popular mangas
// @Description Retrieve a list of popular mangas
// @Tags Manga
// @ID get-popular-manga
// @Accept  json
// @Produce  json
// @Success 200 {array} models.MangaResp
// @Router /manga/popular [get]
func (m *MangaHandler) Popular(w http.ResponseWriter, r *http.Request) {
	mangas, err := m.serv.ListPopularMangas(r.Context())
	if err != nil {
		utils.WriteError(w, 500, "Get Popular Mangas")
		return
	}

	resp := models.MangasRespFromDB(mangas)
	utils.WriteJSON(w, 200, &resp)
}

func (m *MangaHandler) Search(w http.ResponseWriter, r *http.Request) {
	op := "handler Search"

	animes, err := m.serv.ListMangas(r.Context())
	if err != nil {
		utils.WriteError(w, 500, op+"LM")
		return
	}

	utils.WriteJSON(w, 200, &animes)
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
// @Param  genres query []string false "Chapter of the Manga"
// @Param  status query string false "Name of the Manga"
// @Param  country query string false "Chapter of the Manga"
// @Param  orderField query string false "field of the Manga"
// @Param  orderSort query string false "sort of the Manga"
// @Param  page query int false "page not 0"
// @Param  perPage query int false "perPage"
// @Success 200 {array} models.MangaResp
// @Router /manga/filter [get]
func (m *MangaHandler) Filter(w http.ResponseWriter, r *http.Request) {
	// op := "Handler Filter"
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("perPage"))

	filter := models.MangaFilter{
		Name:       q.Get("name"),
		Genres:     q["genres[]"],
		Status:     q.Get("status"),
		Country:    q.Get("country"),
		OrderField: q.Get("orderField"),
		OrderSort:  q.Get("orderSort"),
		Page:       page,
		PerPage:    perPage,
	}

	mangas, err := m.serv.FilterMangas(r.Context(), filter)
	if err != nil {
		utils.WriteError(w, 500, "Filter Mangas")
		return
	}

	utils.WriteJSON(w, 200, mangas)
}
