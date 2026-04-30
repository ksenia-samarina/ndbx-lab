package logout

import (
	"errors"
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
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("No cookie: %v", err)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Invalid cookie: %v", err)
		return
	}
	sid := model.NewSid(cookie.Value)

	exists, _ := h.domain.GetSession(ctx, sid)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = h.domain.DeleteSession(ctx, sid)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error delete auth: %v", err)
		return
	}
	h.writeSessionResponse(w, cookie.Value, 0, http.StatusNoContent)
	return
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
