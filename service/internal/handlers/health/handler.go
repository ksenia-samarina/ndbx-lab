package health

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const OkHealthCheck StatusHealthCheck = "ok"

type Handler struct {
	ttl time.Duration
}

func New(ttl time.Duration) *Handler {
	return &Handler{
		ttl: ttl,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("X-Session-Id")
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "X-Session-Id",
			Value:    cookie.Value,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(h.ttl.Seconds()),
		})
	}
	w.WriteHeader(http.StatusOK)

	healthResp := &Resp{
		Status: OkHealthCheck,
	}

	err = json.NewEncoder(w).Encode(healthResp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
}
