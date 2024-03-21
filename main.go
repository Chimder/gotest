package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return
	}

	mux := http.NewServeMux()

	ctx := context.Background()
	dbl := os.Getenv("DB_URL")

	fmt.Println(dbl)
	conn, err := pgx.Connect(ctx, dbl)
	log.Printf("db conn", conn)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return
	}
	defer conn.Close(ctx)

	// handler := &handler.Manga{}

	// mux.HandleFunc("GET /", handler.Allmanga)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		query := `SELECT id FROM Anime `
		rows, err := conn.Query(ctx, query)
		if err != nil {

			// return nil, fmt.Errorf("unable to query users: %w", err)
		}
		log.Println(rows)
		defer rows.Close()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Listening...")
	http.ListenAndServe(":5000", mux)
}
