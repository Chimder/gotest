package handler

import (
	"net/http"
	"time"

	"github.com/chimas/GoProject/internal/service"
	"github.com/chimas/GoProject/utils"
	"github.com/lib/pq"
)

type MangaHandler struct {
	serv *service.MangaService
}

func NewMangaHandler(s *service.MangaService) *MangaHandler {
	return &MangaHandler{serv: s}
}

type Manga struct {
	Name          string         `json:"name"`
	Img           string         `json:"img"`
	ImgHeader     string         `json:"imgHeader" db:"imgHeader"`
	Describe      string         `json:"describe"`
	Genres        pq.StringArray `json:"genres" db:"genres"`
	Author        string         `json:"author"`
	Country       string         `json:"country"`
	Published     int            `json:"published"`
	AverageRating float64        `json:"averageRating" db:"averageRating"`
	RatingCount   int            `json:"ratingCount" db:"ratingCount"`
	Status        string         `json:"status"`
	Popularity    int            `json:"popularity"`
	Id            int            `json:"id"`
	Chapters      []Chapter      `json:"chapters"`
}

type Chapter struct {
	Chapter   int            `json:"chapter"`
	Img       pq.StringArray `json:"img" db:"img"`
	Name      string         `json:"name"`
	AnimeName string         `json:"animeName" db:"animeName"`
	CreatedAt time.Time      `json:"createdAt" db:"createdAt"`
}

// @Summary Get all mangas
// @Description Retrieve a list of all mangas
// @Tags Manga
// @ID get-all-mangas
// @Accept  json
// @Produce  json
// @Success 200 {array} MangaSwag
// @Router /manga/many [get]
func (m *MangaHandler) Mangas(w http.ResponseWriter, r *http.Request) {
	op := "handler Mangas"

	mangas, err := m.serv.ListMangas(r.Context())
	if err != nil {
		utils.WriteError(w, 500, op+"LM", err)
		return
	}
	if err := utils.WriteJSON(w, 200, &mangas); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
		return
	}
}

// @Summary Get a manga by name
// @Description Retrieve a manga by its name
// @Tags Manga
// @ID get-manga-by-name
// @Accept  json
// @Produce  json
// @Param  name query string true "Name of the Manga"
// @Success 200 {object} MangaSwag
// @Router /manga [get]
func (m *MangaHandler) Manga(w http.ResponseWriter, r *http.Request) {
	op := "handler Manga"
	name := r.URL.Query().Get("name")

	manga, err := m.serv.GetMangaByName(r.Context(), name)
	if err != nil {
		utils.WriteError(w, 500, op+"GMBN", err)
		return
	}

	if err := utils.WriteJSON(w, 200, &manga); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
		return
	}
}

// @Summary Get popular mangas
// @Description Retrieve a list of popular mangas
// @Tags Manga
// @ID get-popular-manga
// @Accept  json
// @Produce  json
// @Success 200 {array} MangaSwag
// @Router /manga/popular [get]
func (m *MangaHandler) Popular(w http.ResponseWriter, r *http.Request) {
	op := "handler Popular"

	animes, err := m.serv.ListPopularMangas(r.Context())
	if err != nil {
		utils.WriteError(w, 500, op+"LPM", err)
		return
	}

	if err := utils.WriteJSON(w, 200, &animes); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
		return
	}
}

func (m *MangaHandler) Search(w http.ResponseWriter, r *http.Request) {
	op := "handler Search"

	animes, err := m.serv.ListMangas(r.Context())
	if err != nil {
		utils.WriteError(w, 500, op+"LM", err)
		return
	}

	if err := utils.WriteJSON(w, 200, &animes); err != nil {
		utils.WriteError(w, 500, op+"WJ", err)
		return
	}
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
// @Success 200 {array} MangaSwag
// @Router /manga/filter [get]
// func (m *MangaHandler) Filter(w http.ResponseWriter, r *http.Request) {
// 	op := "Handler Filter"
// 	params := r.URL.Query()
// 	name := params.Get("name")
// 	genres := params["genres[]"]
// 	status := params.Get("status")
// 	country := params.Get("country")
// 	orderField := params.Get("orderField")
// 	orderSort := params.Get("orderSort")

// 	page, err := strconv.Atoi(params.Get("page"))
// 	if err != nil {
// 		log.Println("not have page")
// 	}
// 	perPage, err := strconv.Atoi(params.Get("perPage"))
// 	if err != nil {
// 		log.Println("not have perPage")
// 	}

// 	var mangas []Manga
// 	query := `SELECT * FROM "Anime"`
// 	args := []interface{}{}
// 	i := 1

// 	if name != "" || status != "" || country != "" || (len(genres) > 0 && genres[0] != "") {
// 		query += " WHERE"
// 	}

// 	if name != "" {
// 		name = "%" + name + "%"
// 		query += fmt.Sprintf(` "name" ILIKE $%d AND`, i)
// 		args = append(args, name)
// 		i++
// 	}
// 	if status != "" {
// 		query += fmt.Sprintf(` "status" = $%d AND`, i)
// 		args = append(args, status)
// 		i++
// 	}
// 	if country != "" {
// 		query += fmt.Sprintf(` "country" = $%d AND`, i)
// 		args = append(args, country)
// 		i++
// 	}
// 	if len(genres) > 0 && genres[0] != "" {
// 		for _, genre := range genres {
// 			query += fmt.Sprintf(` "genres" @> ARRAY[$%d] AND`, i)
// 			args = append(args, genre)
// 			i++
// 		}
// 	}

// 	query = strings.TrimSuffix(query, "AND")
// 	if orderField != "" && orderSort != "" {
// 		query += fmt.Sprintf(` ORDER BY "%s" %s`, orderField, orderSort)
// 	}
// 	if page > 0 && perPage > 0 {
// 		query += fmt.Sprintf(` LIMIT %d OFFSET %d`, perPage, (page-1)*perPage)
// 	}

// 	err = m.sqlx.Select(&mangas, query, args...)
// 	if err != nil {
// 		utils.WriteError(w, 500, op+"SEL", err)
// 		return
// 	}

// 	if err := utils.WriteJSON(w, 200, &mangas); err != nil {
// 		utils.WriteError(w, 500, op+"WJ", err)
// 		return
// 	}
// }
