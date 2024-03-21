package handler

import (
	"encoding/json"
	"net/http"
)

type Manga struct {
	Name          string `json:"name"`
	// Img           string
	// ImgHeader     string
	// Describe      string
	// Genres        []string
	// Author        string
	// Country       string
	// Published     int
	// AverageRating float32
	// RatingCount   int
	// Status        string
	// Popularity    int
	// Id            int
}



func (m *Manga) Allmanga(w http.ResponseWriter, r *http.Request) {

	msg := Manga{Name: "HI"}



	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
