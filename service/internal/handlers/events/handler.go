package events

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"samarina/ndbx/internal/domains/types"
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("X-Session-Id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("No cookie: %v", err)
		return
	}
	sid := types.NewSid(cookie.Value)

	if r.Method == http.MethodPost {
		var event types.Event
		err = json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Invalid JSON: %v", err)
			return
		}

		if event.Address != "" && event.Location.Address == "" {
			event.Location.Address = event.Address
		}
		if event.Location.Address != "" && event.Address == "" {
			event.Address = event.Location.Address
		}

		// 400
		if event.Title == "" {
			h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			eventRespMsg := &Resp{
				Message: StatusEventCreation(fmt.Sprintf(string(invalidFieldName), "title")),
			}
			err = json.NewEncoder(w).Encode(eventRespMsg)
			return
		}
		if event.Location.Address == "" {
			h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			eventRespMsg := &Resp{
				Message: StatusEventCreation(fmt.Sprintf(string(invalidFieldName), "address")),
			}
			err = json.NewEncoder(w).Encode(eventRespMsg)
			return
		}
		if event.StartedAt == "" {
			h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			eventRespMsg := &Resp{
				Message: StatusEventCreation(fmt.Sprintf(string(invalidFieldName), "started_at")),
			}
			err = json.NewEncoder(w).Encode(eventRespMsg)
			return
		}
		if event.FinishedAt == "" {
			h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			eventRespMsg := &Resp{
				Message: StatusEventCreation(fmt.Sprintf(string(invalidFieldName), "finished_at")),
			}
			err = json.NewEncoder(w).Encode(eventRespMsg)
			return
		}

		// 401
		exists, _ := h.domain.GetSession(ctx, sid)
		if !exists {
			h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}
		userID, err := h.domain.GetInternalUserID(ctx, sid)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Invalid user ID: %v", err)
			return
		}
		if userID == "" {
			h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
			return
		}

		// 409
		createdBy := userID
		eventDB, err := h.domain.GetEvent(ctx, createdBy)
		if eventDB != nil && eventDB.Title == event.Title {
			h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
			eventRespMsg := &Resp{
				Message: eventAlreadyExists,
			}
			err = json.NewEncoder(w).Encode(eventRespMsg)
			return
		}

		// 201
		eventID, err := h.domain.CreateEvent(ctx, createdBy, &event)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Invalid create event DB: %v", err)
			return
		}
		err = h.domain.UpdateUserSession(ctx, sid, h.ttl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error updating user session: %v", err)
			return
		}
		h.writeSessionResponse(w, cookie.Value, h.ttl, http.StatusCreated)
		eventCreationMsg := &EventID{
			ID: eventID,
		}
		err = json.NewEncoder(w).Encode(eventCreationMsg)
		return
	}
	if r.Method == http.MethodGet {
		query := r.URL.Query()

		title := query.Get("title")

		limitStr := query.Get("limit")
		limit := int64(-1)
		if limitStr != "" {
			l, err := strconv.ParseInt(limitStr, 10, 64)
			if err != nil || l < 0 {
				h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
				eventRespMsg := &Resp{
					Message: StatusEventCreation(fmt.Sprintf(string(invalidParameterName), "limit")),
				}
				err = json.NewEncoder(w).Encode(eventRespMsg)
				return
			}
			limit = l
		}

		offsetStr := query.Get("offset")
		offset := int64(0)
		if offsetStr != "" {
			off, err := strconv.ParseInt(offsetStr, 10, 64)
			if err != nil || off < 0 {
				h.writeSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
				eventRespMsg := &Resp{
					Message: StatusEventCreation(fmt.Sprintf(string(invalidParameterName), "offset")),
				}
				err = json.NewEncoder(w).Encode(eventRespMsg)
				return
			}
			offset = off
		}

		events, err := h.domain.ListEvents(ctx, title, offset, limit)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error listing events: %v", err)
			return
		}

		resp := map[string]interface{}{
			"events": events,
			"count":  len(events),
		}
		h.writeSessionResponse(w, cookie.Value, h.ttl, http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		return
	}
}

func (h *Handler) writeSessionResponse(w http.ResponseWriter, sid string, ttl time.Duration, status int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "X-Session-Id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(ttl.Seconds()),
	})
	if status != http.StatusUnauthorized {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)
}
