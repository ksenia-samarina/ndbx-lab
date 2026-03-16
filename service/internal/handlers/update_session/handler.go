package update_session

import (
	"log"
	"net/http"
	"os"
	"samarina/ndbx/internal/domains/session"
	"time"
)

type Handler struct {
	Domain Domain
}

func New(domain Domain) *Handler {
	return &Handler{
		Domain: domain,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Записываем общие заголовки
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	ttl, _ := time.ParseDuration(os.Getenv("APP_USER_SESSION_TTL") + "s")

	cookie, err := r.Cookie("X-Session-Id")
	// Куки представлена, записываем новую куки
	if err == nil {
		sid := session.Id{HexString: cookie.Value}
		exists, _ := h.Domain.CheckSessionExist(ctx, sid)
		// Сессия уже существует, обновляем сессию
		if exists {
			err := h.Domain.UpdateSession(ctx, sid, ttl)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error update session: %v", err)
				return
			}
			h.writeSessionResponse(w, sid.HexString, ttl, http.StatusOK)
			return
		}
	}
	// Куки не представлена, либо сессия не существует -> создаем новую сессию
	newSession, err := h.Domain.CreateSession(ctx, ttl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error create session: %v", err)
		return
	}
	h.writeSessionResponse(w, newSession.HexString, ttl, http.StatusCreated)
}

func (h *Handler) writeSessionResponse(w http.ResponseWriter, sid string, ttl time.Duration, status int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "X-Session-Id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(ttl.Seconds()),
	})
	w.WriteHeader(status)
}
