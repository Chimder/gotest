package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func NewMangaHandler(db *sqlx.DB, rdb *redis.Client) *MangaHandler {
	return &MangaHandler{db: db, rdb: rdb}
}

type MangaHandler struct {
	db  *sqlx.DB
	rdb *redis.Client
}

type Manga struct {
	Name          string         `json:"name"`
	Img           string         `json:"img"`
	ImgHeader     string         `json:"imgHeader" db:"imgHeader"`
	Describe      string         `json:"describe"`
	Genres        pq.StringArray `json:"genres"`
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
	Img       pq.StringArray `json:"img"`
	Name      string         `json:"name"`
	AnimeName string         `json:"animeName" db:"animeName"`
	CreatedAt time.Time      `json:"created" db:"createdAt"`
}

func (m *MangaHandler) Mangas(w http.ResponseWriter, r *http.Request) {

	query := `SELECT * FROM "Anime"`
	var mangas []Manga
	err := m.db.Select(&mangas, query)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mangas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (m *MangaHandler) Manga(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	name := r.PathValue("name")

	val, err := m.rdb.Get(ctx, name).Result()
	if err == redis.Nil {

		query := `SELECT * FROM "Anime" WHERE name=$1`
		chaptersQuery := `SELECT * FROM "Chapter" WHERE "animeName" =$1`
		var manga Manga
		err := m.db.Get(&manga, query, name)
		if err != nil {
			log.Fatal(err)
		}
		var chapters []Chapter
		err = m.db.Select(&chapters, chaptersQuery, name)
		if err != nil {
			log.Fatal(err)
		}

		manga.Chapters = chapters

		mangaJSON, err := json.Marshal(manga)
		if err != nil {
			log.Fatal(err)
		}

		err = m.rdb.Set(ctx, name, mangaJSON, time.Minute).Err()
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(manga); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		manga := Manga{}
		err := json.Unmarshal([]byte(val), &manga)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(manga); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}


func (m *MangaHandler) Chapter(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	chapt := r.PathValue("chapter")

	var chapter Chapter

	query := `SELECT * FROM "Chapter" WHERE "animeName" =$1 AND chapter=$2`
	
	err := m.db.Get(&chapter, query, name, chapt)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chapter); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (m *MangaHandler) Popular(w http.ResponseWriter, r *http.Request) {

	query := `SELECT * FROM "Anime"`
	var animes []Manga
	err := m.db.Select(&animes, query)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (m *MangaHandler) Search(w http.ResponseWriter, r *http.Request) {

	query := `SELECT * FROM "Anime"`
	var animes []Manga
	err := m.db.Select(&animes, query)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (m *MangaHandler) Rating(w http.ResponseWriter, r *http.Request) {

	query := `SELECT * FROM "Anime"`
	var animes []Manga
	err := m.db.Select(&animes, query)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
