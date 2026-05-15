package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"samarina/ndbx/internal/domains/validator"
	"samarina/ndbx/internal/handlers/utils"
	"samarina/ndbx/internal/model"
	"time"
)

type Handler struct {
	eventDomain     eventDomain
	validatorDomain validatorDomain
	ttl             time.Duration
}

func New(eventDomain eventDomain, validatorDomain validatorDomain, ttl time.Duration) *Handler {
	return &Handler{
		eventDomain:     eventDomain,
		validatorDomain: validatorDomain,
		ttl:             ttl,
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

		var target *validator.ErrInvalidFieldName
		filter, err := h.validatorDomain.ValidateParams(query)
		if errors.As(err, &target) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: target.Field})
		}

		var createdBy = ""
		username := query.Get("user")
		if username != "" {
			user, err := h.eventDomain.GetUserByUsername(ctx, username)
			if err == nil {
				createdBy = user.ID.Hex()
			}
		}
		filter.User = createdBy

		eventsList, _ := h.eventDomain.GetEvents(ctx, filter)

		resp := map[string]interface{}{
			"events": eventsList,
			"count":  len(eventsList),
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

		userID, err := h.eventDomain.GetInternalUserID(ctx, sid)
		if err != nil { // TODO: кастомная ошибка
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}

		events, err := h.eventDomain.GetEventsByUserID(ctx, userID)
		for _, e := range events {
			if e.Title == event.Title {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
				utils.EncodeErrorResponse(w, ErrEventAlreadyExists)
				return
			}
		}

		eventID, _ := h.eventDomain.RegisterEvent(ctx, userID, event)
		err = h.eventDomain.UpdateUserSession(ctx, sid, h.ttl)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}
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

		validCategories := map[string]struct{}{"": {}, "meetup": {}, "concert": {}, "exhibition": {}, "party": {}, "other": {}}
		if _, exists := validCategories[event.Category]; !exists {
			err := h.eventDomain.UpdateUserSession(ctx, sid, h.ttl)
			if err != nil {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
				return
			}
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "category"})
			return
		}

		if event.Price < 0 {
			err := h.eventDomain.UpdateUserSession(ctx, sid, h.ttl)
			if err != nil {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
				return
			}
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "price"})
			return
		}

		existedEvent, err := h.eventDomain.GetEventByID(ctx, id)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			utils.EncodeErrorResponse(w, ErrEventNotFound)
			return
		}

		currentUserID, err := h.eventDomain.GetInternalUserID(ctx, sid)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}

		if existedEvent.CreatedBy != currentUserID {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			utils.EncodeErrorResponse(w, ErrEventNotFound)
			return
		}

		err = h.eventDomain.UpdateEventsLocationCity(ctx, id, event.Location.City)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNoContent)
	case http.MethodGet:
		event, err := h.eventDomain.GetEventByID(ctx, id)
		if err != nil { // TODO: custom error
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			utils.EncodeErrorResponse(w, ErrEventNotExist)
			return
		}

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
		_ = json.NewEncoder(w).Encode(event)
	}
}
