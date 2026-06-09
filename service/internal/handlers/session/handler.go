package session

import (
	"log"
	"net/http"
	"samarina/ndbx/internal/model"
	"time"
)

type Handler struct {
	domain domain
	ttl    time.Duration
}

func New(domain domain, ttl time.Duration) *Handler {
	return &Handler{
		domain: domain,
		ttl:    ttl,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("X-Session-Id")
	if err == nil {
		sid := model.NewSid(cookie.Value)
		exists, _ := h.domain.GetSession(ctx, sid)
		if exists {
			err := h.domain.UpdateSession(ctx, sid, h.ttl)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error update auth: %v", err)
				return
			}
			h.writeSessionResponse(w, cookie.Value, h.ttl, http.StatusOK)
			return
		}
	}
	newSession, err := h.domain.CreateSession(ctx, h.ttl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error create auth: %v", err)
		return
	}
	h.writeSessionResponse(w, newSession.HexString, h.ttl, http.StatusCreated)
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
