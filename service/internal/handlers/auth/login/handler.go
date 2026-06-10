package login

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"samarina/ndbx/internal/model"
	"samarina/ndbx/internal/utils"
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
	var cookieValue string
	if cookie != nil {
		cookieValue = cookie.Value
	} else {
		cookieValue = ""
	}
	sid := model.NewSid(cookieValue)

	var login model.Login
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Invalid JSON: %v", err)
		return
	}

	// 400
	if login.Username == "" {
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		userRegistrationResp := &Resp{
			Message: StatusUserLogin(fmt.Sprintf(string(invalidFieldName), "username")),
		}
		err = json.NewEncoder(w).Encode(userRegistrationResp)
		return
	}
	if login.Password == "" {
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		userRegistrationResp := &Resp{
			Message: StatusUserLogin(fmt.Sprintf(string(invalidFieldName), "password")),
		}
		err = json.NewEncoder(w).Encode(userRegistrationResp)
		return
	}

	// 401
	user, err := h.domain.GetByUsername(ctx, login.Username)
	if err != nil {
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
		loginResp := &Resp{
			Message: userInvalidCredentials,
		}
		err = json.NewEncoder(w).Encode(loginResp)
		return
	}
	if !utils.CheckPasswordHash(login.Password, user.PasswordHash) {
		h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
		loginResp := &Resp{
			Message: userInvalidCredentials,
		}
		err = json.NewEncoder(w).Encode(loginResp)
		return
	}
	// 204
	userID, err := h.domain.GetInternalUserID(ctx, user.Username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error get user id: %v", err)
		return
	}
	exists, _ := h.domain.GetSession(ctx, sid)
	if exists {
		err := h.domain.UpdateUserSession(ctx, userID, sid, h.ttl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error update auth: %v", err)
			return
		}
		h.writeSessionResponse(w, cookieValue, h.ttl, http.StatusNoContent)
		return
	}

	newSid, err := h.domain.CreateUserSession(ctx, userID, h.ttl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error create auth: %v", err)
		return
	}

	h.writeSessionResponse(w, newSid.HexString, h.ttl, http.StatusNoContent)
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
