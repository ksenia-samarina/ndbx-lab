package reactions

import (
	"encoding/json"
	"errors"
	"net/http"
	"samarina/ndbx/internal/domains/reviews"
	"samarina/ndbx/internal/handlers/utils"
	"samarina/ndbx/internal/model"
	"samarina/ndbx/internal/repository/cassandra"
	"strconv"
	"time"
)

type Handler struct {
	eventDomain  eventDomain
	reviewDomain reviewDomain
	ttl          time.Duration
}

func New(eventDomain eventDomain, reviewDomain reviewDomain, ttl time.Duration) *Handler {
	return &Handler{
		eventDomain:  eventDomain,
		reviewDomain: reviewDomain,
		ttl:          ttl,
	}
}

type CreateReviewRequest struct {
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
}

type GetReviewsResponse struct {
	Reviews []model.Review `json:"reviews"`
	Count   int            `json:"count"`
}

func (h *Handler) PostReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("X-Session-Id")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sid := model.NewSid(cookie.Value)
	createdBy, err := h.eventDomain.GetInternalUserID(ctx, sid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid fields"})
		return
	}

	eventID := r.PathValue("event_id")

	reviewID, err := h.reviewDomain.AddReview(ctx, eventID, createdBy, req.Comment, req.Rating)
	var invalidFieldErr *reviews.ErrInvalidFieldName
	if err != nil {
		if errors.As(err, &invalidFieldErr) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			_ = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: invalidFieldErr.Field})
			return
		}
		if errors.Is(err, reviews.ErrEventNotFound) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			_ = utils.EncodeErrorResponse(w, ErrEventNotFound)
			return
		}
		if errors.Is(err, cassandra.ErrAlreadyExists) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
			_ = utils.EncodeErrorResponse(w, ErrAlreadyExists)
			return
		}
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
		return
	}

	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": reviewID})
}

func (h *Handler) GetReviews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var sidStr string
	if cookie, err := r.Cookie("X-Session-Id"); err == nil {
		sidStr = cookie.Value
	}

	q := r.URL.Query()
	limitStr := q.Get("limit")
	offsetStr := q.Get("offset")

	limit := 10
	offset := 0
	var err error

	if limitStr != "" {
		if limit, err = strconv.Atoi(limitStr); err != nil || limit < 0 {
			utils.WriteSessionResponse(w, sidStr, h.ttl, http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": `invalid "limit" field`})
			return
		}
	}
	if offsetStr != "" {
		if offset, err = strconv.Atoi(offsetStr); err != nil || offset < 0 {
			utils.WriteSessionResponse(w, sidStr, h.ttl, http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": `invalid "offset" field`})
			return
		}
	}

	eventID := r.PathValue("event_id")

	reviewsList, count, err := h.reviewDomain.GetReviews(ctx, eventID, limit, offset)
	var invalidFieldErr *reviews.ErrInvalidFieldName
	if err != nil {
		if errors.As(err, &invalidFieldErr) {
			utils.WriteSessionResponse(w, sidStr, h.ttl, http.StatusBadRequest)
			_ = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: invalidFieldErr.Field})
			return
		}
		if errors.Is(err, reviews.ErrEventNotFound) {
			utils.WriteSessionResponse(w, sidStr, h.ttl, http.StatusNotFound)
			_ = utils.EncodeErrorResponse(w, ErrEventNotFound)
			return
		}
		if errors.Is(err, cassandra.ErrAlreadyExists) {
			utils.WriteSessionResponse(w, sidStr, h.ttl, http.StatusConflict)
			_ = utils.EncodeErrorResponse(w, ErrAlreadyExists)
			return
		}
		utils.WriteSessionResponse(w, sidStr, h.ttl, http.StatusInternalServerError)
		return
	}

	if reviewsList == nil {
		reviewsList = []model.Review{}
	}

	utils.WriteSessionResponse(w, sidStr, h.ttl, http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetReviewsResponse{
		Reviews: reviewsList,
		Count:   count,
	})
}

type PatchReviewRequest struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment"`
}

func (h *Handler) PatchReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("X-Session-Id")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sid := model.NewSid(cookie.Value)
	createdBy, err := h.eventDomain.GetInternalUserID(ctx, sid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req PatchReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		return
	}

	eventID := r.PathValue("event_id")
	reviewID := r.PathValue("review_id")

	err = h.reviewDomain.PatchReview(ctx, eventID, reviewID, createdBy, req.Rating, req.Comment)
	if err != nil {
		var invalidFieldErr *reviews.ErrInvalidFieldName
		if errors.As(err, &invalidFieldErr) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			_ = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: invalidFieldErr.Field})
			return
		}
		if errors.Is(err, reviews.ErrEventNotFound) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
			_ = utils.EncodeErrorResponse(w, ErrEventNotFound)
			return
		}
		if errors.Is(err, cassandra.ErrAlreadyExists) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
			_ = utils.EncodeErrorResponse(w, ErrAlreadyExists)
			return
		}
		if errors.Is(err, reviews.ErrNoReviewAccess) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
		return
	}

	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNoContent)
}
