package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"samarina/ndbx/internal/domains/types"
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

	cookie, _ := r.Cookie("X-Session-Id")
	var sid *types.Sid
	if cookie != nil {
		sid = types.NewSid(cookie.Value)
	} else {
		sid = types.NewSid("")
	}

	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Invalid JSON: %v", err)
		return
	}

	// 400
	if user.Username == "" {
		err := h.domain.UpdateUserSession(ctx, sid, h.ttl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error updating user session: %v", err)
			return
		}
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		userRegistrationResp := &Resp{
			Message: StatusUserRegistration(fmt.Sprintf(string(invalidFieldName), "username")),
		}
		err = json.NewEncoder(w).Encode(userRegistrationResp)
		return
	}
	if user.FullName == "" {
		err := h.domain.UpdateUserSession(ctx, sid, h.ttl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error updating user session: %v", err)
			return
		}
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		userRegistrationResp := &Resp{
			Message: StatusUserRegistration(fmt.Sprintf(string(invalidFieldName), "full_name")),
		}
		err = json.NewEncoder(w).Encode(userRegistrationResp)
		return
	}
	if user.Password == "" {
		err := h.domain.UpdateUserSession(ctx, sid, h.ttl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error updating user session: %v", err)
			return
		}
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		userRegistrationResp := &Resp{
			Message: StatusUserRegistration(fmt.Sprintf(string(invalidFieldName), "password")),
		}
		err = json.NewEncoder(w).Encode(userRegistrationResp)
		return
	}

	existingUser, err := h.domain.GetByUsername(ctx, user.Username)
	// 409
	if err == nil && existingUser != nil {
		err = h.domain.UpdateUserSession(ctx, sid, h.ttl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error update session: %v", err)
			return
		}
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
		userRegistrationResp := &Resp{
			Message: userAlreadyExists,
		}
		err = json.NewEncoder(w).Encode(userRegistrationResp)
		return
	}
	// 201
	newSession, err := h.domain.RegisterUser(ctx, &user, h.ttl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error create session: %v", err)
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
