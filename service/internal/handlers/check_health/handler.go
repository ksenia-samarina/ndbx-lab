package check_health

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"samarina/ndbx/internal/dto/response"
	"time"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Записываем все заголовки
	w.Header().Set("Content-Type", "application/json")

	// Логика обработки куки
	cookie, err := r.Cookie("X-Session-Id")
	// Куки представлена, записываем новую куки
	if err == nil {
		ttl, err := time.ParseDuration(os.Getenv("APP_USER_SESSION_TTL") + "s")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error parsing APP_USER_SESSION_TTL: %v", err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "X-Session-Id",
			Value:    cookie.Value,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(ttl.Seconds()),
		})
	}
	w.WriteHeader(http.StatusOK)
	// Записываем тело ответа
	healthCheck := response.HealthCheck{
		Status: response.OkHealthCheck,
	}
	err = json.NewEncoder(w).Encode(healthCheck)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
}
