package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"samarina/ndbx/internal/domains/validator"
	"samarina/ndbx/internal/handlers/utils"
	"samarina/ndbx/internal/model"
	"strings"
	"time"
)

type reviewsDomain interface {
	EnrichEventsWithReviews(ctx context.Context, events []model.Event, includeReviews bool) ([]model.Event, error)
}

type Handler struct {
	eventDomain     eventDomain
	reactionsDomain reactionsDomain
	reviewsDomain   reviewsDomain
	validatorDomain validatorDomain
	ttl             time.Duration
}

func New(eventDomain eventDomain, reactionsDomain reactionsDomain, reviewsDomain reviewsDomain, validatorDomain validatorDomain, ttl time.Duration) *Handler {
	return &Handler{
		eventDomain:     eventDomain,
		reactionsDomain: reactionsDomain,
		reviewsDomain:   reviewsDomain,
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

		includeParam := query.Get("include")
		includeReactions := strings.Contains(includeParam, "reactions")
		includeReviews := strings.Contains(includeParam, "reviews")
		query.Del("include")

		var target *validator.ErrInvalidFieldName
		filter, err := h.validatorDomain.ValidateParams(query)
		if errors.As(err, &target) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: target.Field})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
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

		enrichedEvents, err := h.reactionsDomain.EnrichEventsWithReactions(ctx, eventsList, includeReactions)
		if err != nil {
			log.Printf("Failed to enrich events with reactions: %v", err)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}

		enrichedEvents, err = h.reviewsDomain.EnrichEventsWithReviews(ctx, enrichedEvents, includeReviews)
		if err != nil {
			log.Printf("Failed to enrich events with reviews: %v", err)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"events": enrichedEvents,
			"count":  len(enrichedEvents),
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
			err := utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "title"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}
		if event.Location.Address == "" {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err := utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "address"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}
		if event.StartedAt == "" {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err := utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "started_at"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}
		if event.FinishedAt == "" {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err := utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "finished_at"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		userID, err := h.eventDomain.GetInternalUserID(ctx, sid)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}

		events, err := h.eventDomain.GetEventsByUserID(ctx, userID)
		for _, e := range events {
			if e.Title == event.Title {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
				err := utils.EncodeErrorResponse(w, ErrEventAlreadyExists)
				if err != nil {
					log.Printf("Encode error wasn't sent to client: %v", err)
				}
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
		err = json.NewEncoder(w).Encode(map[string]string{
			"id": fmt.Sprintf("%v", eventID),
		})
		if err != nil {
			log.Printf("Encode error wasn't sent to client: %v", err)
		}
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
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "category"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		if event.Price < 0 {
			err := h.eventDomain.UpdateUserSession(ctx, sid, h.ttl)
			if err != nil {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
				return
			}
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "price"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		existedEvent, err := h.eventDomain.GetEventByID(ctx, id)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			err = utils.EncodeErrorResponse(w, ErrEventNotFound)
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		currentUserID, err := h.eventDomain.GetInternalUserID(ctx, sid)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}

		if existedEvent.CreatedBy != currentUserID {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			err = utils.EncodeErrorResponse(w, ErrEventNotFound)
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		err = h.eventDomain.UpdateEventsLocationCity(ctx, id, event.Location.City)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNoContent)
	case http.MethodGet:
		includeParam := r.URL.Query().Get("include")
		includeReactions := strings.Contains(includeParam, "reactions")
		includeReviews := strings.Contains(includeParam, "reviews")

		event, err := h.eventDomain.GetEventByID(ctx, id)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			err = utils.EncodeErrorResponse(w, ErrEventNotFound)
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		enrichedEvents, err := h.reactionsDomain.EnrichEventsWithReactions(ctx, []model.Event{event}, includeReactions)
		if err != nil {
			log.Printf("Failed to enrich single event with reactions: %v", err)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}

		enrichedEvents, err = h.reviewsDomain.EnrichEventsWithReviews(ctx, enrichedEvents, includeReviews)
		if err != nil {
			log.Printf("Failed to enrich single event with reviews: %v", err)
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}

		if len(enrichedEvents) == 0 {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			_ = utils.EncodeErrorResponse(w, ErrEventNotFound)
			return
		}

		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
		_ = json.NewEncoder(w).Encode(enrichedEvents[0])
	}
}
