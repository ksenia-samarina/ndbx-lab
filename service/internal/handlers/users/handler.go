package users

import (
	"encoding/json"
	"net/http"
	"samarina/ndbx/internal/handlers/utils"
	"samarina/ndbx/internal/model"
	"strconv"
	"time"
)

const DateLayout = "20060102"

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
	var hexString string
	if cookie != nil {
		hexString = cookie.Value
	}
	sid := model.NewSid(hexString)

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
		if err == nil { // TODO: кастомная ошибка что нет пользака
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
	var hexString string
	if cookie != nil {
		hexString = cookie.Value
	}
	sid := model.NewSid(hexString)

	id := r.PathValue("id")
	user, err := h.domain.GetUserByUserID(ctx, id)
	if err != nil { // TODO: кастомная ошибка что нет пользака
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		utils.EncodeErrorResponse(w, ErrUserNotFound)
		return
	}
	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
	_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
}

func (h *Handler) GetUserEventsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, _ := r.Cookie("X-Session-Id")
	var hexString string
	if cookie != nil {
		hexString = cookie.Value
	}
	sid := model.NewSid(hexString)

	id := r.PathValue("id")
	_, err := h.domain.GetUserByUserID(ctx, id)
	if err != nil { // TODO: кастомная ошибка что нет пользака
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		utils.EncodeErrorResponse(w, ErrUserNotFound)
		return
	}

	query := r.URL.Query()

	limitStr := query.Get("limit")
	limit := int64(-1)
	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || l < 0 {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "limit"})
			return
		}
		limit = l
	}

	offsetStr := query.Get("offset")
	offset := int64(0)
	if offsetStr != "" {
		off, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || off < 0 {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "offset"})
			return
		}
		offset = off
	}

	category := query.Get("category")
	validCategories := map[string]bool{"": true, "meetup": true, "concert": true, "exhibition": true, "party": true, "other": true}
	if !validCategories[category] {
		_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "category"})
		return
	}

	var priceFrom int64 = -1
	if pFrom := query.Get("price_from"); pFrom != "" {
		val, err := strconv.ParseInt(pFrom, 10, 64)
		if err != nil || val < 0 {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "price_from"})
			return
		}
		priceFrom = val
	}

	var priceTo int64 = -1
	if pTo := query.Get("price_to"); pTo != "" {
		val, err := strconv.ParseInt(pTo, 10, 64)
		if err != nil || val < 0 {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "price_to"})
			return
		}
		priceTo = val
	}

	var dateFrom string
	if dFrom := query.Get("date_from"); dFrom != "" {
		t, err := time.Parse(DateLayout, dFrom)
		if err != nil {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "date_from"})
			return
		}
		dateFrom = t.Format(time.RFC3339)
	}

	var dateTo string
	if dTo := query.Get("date_to"); dTo != "" {
		t, err := time.Parse(DateLayout, dTo)
		if err != nil {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "date_to"})
			return
		}
		dateTo = t.Format(time.RFC3339)
	}

	filter := model.EventFilter{
		ID:        query.Get("id"),
		Title:     query.Get("title"),
		Category:  query.Get("category"),
		PriceFrom: priceFrom,
		PriceTo:   priceTo,
		City:      query.Get("city"),
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		User:      id,
		Offset:    offset,
		Limit:     limit,
	}

	events, _ := h.domain.GetEvents(ctx, filter)

	resp := map[string]interface{}{
		"events": events,
		"count":  len(events),
	}
	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
	_ = json.NewEncoder(w).Encode(&resp)
}
