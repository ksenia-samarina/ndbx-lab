package users

import (
	"encoding/json"
	"net/http"
	"samarina/ndbx/internal/handlers/utils"
	"samarina/ndbx/internal/model"
	"strconv"
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

func (h *Handler) RegisterOrGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, _ := r.Cookie("X-Session-Id")
	sid := model.NewSid(cookie.Value)

	switch r.Method {
	case http.MethodGet: // Возвращает список организаторов, отвечающий параметрам поиска
		query := r.URL.Query()
		limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if query.Get("limit") != "" && (err != nil || limit < 0) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "limit"})
			return
		}
		if query.Get("limit") == "" {
			limit = 10
		}

		offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if query.Get("offset") != "" && (err != nil || offset < 0) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "offset"})
			return
		}

		name := query.Get("name")
		id := query.Get("id")

		users, _ := h.domain.GetUsers(ctx, id, name, uint64(limit), uint64(offset))

		resp := map[string]interface{}{
			"users": users,
			"count": len(users),
		}

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	case http.MethodPost: // Регистрирует нового пользователя
		var user model.User
		_ = json.NewDecoder(r.Body).Decode(&user)

		if user.Username == "" {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "username"})
			return
		}
		if user.FullName == "" {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "full_name"})
			return
		}
		if user.Password == "" {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "password"})
			return
		}

		_, err := h.domain.GetUserByUsername(ctx, user.Username)
		if err != nil { // TODO: кастомная ошибка что нет пользака
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
			utils.EncodeErrorResponse(w, ErrUserAlreadyExists)
			return
		}

		newSid, _ := h.domain.RegisterUser(ctx, user, h.ttl)
		utils.WriteSessionResponse(w, newSid.HexString, h.ttl, http.StatusCreated)
	}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, _ := r.Cookie("X-Session-Id")
	sid := model.NewSid(cookie.Value)

	id := r.PathValue("id")
	user, err := h.domain.GetUserByUserID(ctx, id)
	if err != nil { // TODO: кастомная ошибка что нет пользака
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		utils.EncodeErrorResponse(w, ErrUserNotFound)
		return
	}
	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUserEventsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, _ := r.Cookie("X-Session-Id")
	sid := model.NewSid(cookie.Value)

	id := r.PathValue("id")
	_, err := h.domain.GetUserByUserID(ctx, id)
	if err != nil { // TODO: кастомная ошибка что нет пользака
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		utils.EncodeErrorResponse(w, ErrUserNotFound)
		return
	}

	events, _ := h.domain.GetUserEventsByUserID(ctx, id)

	resp := map[string]interface{}{
		"events": events,
		"count":  len(events),
	}
	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
	_ = json.NewEncoder(w).Encode(&resp)
}
