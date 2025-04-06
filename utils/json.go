package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "JSON encoding error: %v"}`, err), http.StatusInternalServerError)
	}
}

func WriteJSONRedis(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func WriteError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	slog.Error("HTTP", "err:", err)
	if err := json.NewEncoder(w).Encode(map[string]string{"err": err}); err != nil {
		http.Error(w, `{"err": err encode error}`, status)
	}
}

func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()

	r.Body = http.MaxBytesReader(nil, r.Body, 1048576)

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	return nil
}
