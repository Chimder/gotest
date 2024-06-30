package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

type Jwt struct {
}

var secretKey = []byte("your_secret_key2")

func Encrypt(payload interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["data"] = payload

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("Ошибка при подписании токена: %v", err)
	}

	return tokenString, nil
}

func Decrypt(tokenString string, v interface{}) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return fmt.Errorf("Ошибка при разборе токена: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := claims["data"].(map[string]interface{})
		bytes, _ := json.Marshal(data)
		json.Unmarshal(bytes, &v)
		return nil
	} else {
		return fmt.Errorf("Недействительный токен")
	}
}

// /////////////////////////////////////////////////////////////////////////
type contextKey string
type User struct {
	Id        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Favorite  pq.StringArray `json:"favorite"`
	CreatedAt time.Time      `json:"createdAt" db:"createdAt"`
}

const userContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("manka_google_user")
		log.Println("COOKIE", cookie)
		if err != nil {
			log.Println("COOKIE", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user User
		err = Decrypt(cookie.Value, &user)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}
