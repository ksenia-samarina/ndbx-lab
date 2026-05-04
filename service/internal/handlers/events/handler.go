package events

import (
	"encoding/json"
	"fmt"
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

func (h *Handler) RegisterOrGetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, _ := r.Cookie("X-Session-Id")
	var hexString string
	if cookie != nil {
		hexString = cookie.Value
	}
	sid := model.NewSid(hexString)

	switch r.Method {
	case http.MethodGet:
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

		var dateFrom time.Time
		if dFrom := query.Get("date_from"); dFrom != "" {
			t, err := time.Parse(DateLayout, dFrom)
			if err != nil {
				_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
				utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "date_from"})
				return
			}
			dateFrom = t
		}

		var dateTo time.Time
		if dTo := query.Get("date_to"); dTo != "" {
			t, err := time.Parse(DateLayout, dTo)
			if err != nil {
				_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
				utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "date_to"})
				return
			}
			dateTo = t
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
			User:      query.Get("user"),
			Offset:    offset,
			Limit:     limit,
		}

		events, _ := h.domain.GetEvents(ctx, filter)

		resp := map[string]interface{}{
			"events": events,
			"count":  len(events),
		}

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	case http.MethodPost:
		var event model.Event
		_ = json.NewDecoder(r.Body).Decode(&event)

		if event.Address != "" && event.Location.Address == "" {
			event.Location.Address = event.Address
		}
		if event.Location.Address != "" && event.Address == "" {
			event.Address = event.Location.Address
		}

		if event.Title == "" {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "title"})
			return
		}
		if event.Location.Address == "" {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "address"})
			return
		}
		if event.StartedAt == "" {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "started_at"})
			return
		}
		if event.FinishedAt == "" {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "finished_at"})
			return
		}

		userID, err := h.domain.GetInternalUserID(ctx, sid)
		if err != nil { // TODO: кастомная ошибка
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}

		events, err := h.domain.GetEventsByUserID(ctx, userID)
		for _, e := range events {
			if e.Title == event.Title {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
				utils.EncodeErrorResponse(w, ErrEventAlreadyExists)
				return
			}
		}

		eventID, _ := h.domain.RegisterEvent(ctx, userID, event)
		_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": fmt.Sprintf("%v", eventID),
		})
	}
}

func (h *Handler) GetOrEditEventData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, _ := r.Cookie("X-Session-Id")
	var hexString string
	if cookie != nil {
		hexString = cookie.Value
	}
	sid := model.NewSid(hexString)

	id := r.PathValue("id")

	switch r.Method {
	case http.MethodPatch:
		var event model.Event
		_ = json.NewDecoder(r.Body).Decode(&event)

		validCategories := map[string]bool{"": true, "meetup": true, "concert": true, "exhibition": true, "party": true, "other": true}
		if !validCategories[event.Category] {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "category"})
			return
		}

		if event.Price < 0 {
			_ = h.domain.UpdateUserSession(ctx, sid, h.ttl)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "price"})
			return
		}

		currentUserID, err := h.domain.GetInternalUserID(ctx, sid)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}

		existedEvent, err := h.domain.GetEventByID(ctx, id)
		if err != nil || existedEvent.CreatedBy != currentUserID {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			utils.EncodeErrorResponse(w, ErrEventNotFound)
			return
		}

		_ = h.domain.UpdateEventsLocationCity(ctx, id, event.Location.City)

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNoContent)
	case http.MethodGet:
		event, err := h.domain.GetEventByID(ctx, id)
		if err != nil { // TODO: custom error
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			utils.EncodeErrorResponse(w, ErrEventNotExist)
			return
		}

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
		_ = json.NewEncoder(w).Encode(event)
	}
}
