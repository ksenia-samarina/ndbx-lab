package reactions

import (
	"errors"
	"log"
	"net/http"
	"samarina/ndbx/internal/handlers/utils"
	"samarina/ndbx/internal/model"
	"time"
)

type Handler struct {
	eventDomain     eventDomain
	reactionsDomain reactionsDomain
	ttl             time.Duration
}

func New(eventDomain eventDomain, reactionsDomain reactionsDomain, ttl time.Duration) *Handler {
	return &Handler{
		eventDomain:     eventDomain,
		reactionsDomain: reactionsDomain,
		ttl:             ttl,
	}
}

func (h *Handler) SetLike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("X-Session-Id")
	if errors.Is(err, http.ErrNoCookie) {
		log.Printf("cookie %s not found", cookie.Value)
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("No cookie: %v", err)
		return
	}

	hexString := cookie.Value
	sid := model.NewSid(hexString)

	createdBy, err := h.eventDomain.GetInternalUserID(ctx, sid)
	if err != nil { // TODO: кастомная ошибка
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
		return
	}

	eventId := r.PathValue("event_id")

	_, err = h.eventDomain.GetEventByID(ctx, eventId)
	if err != nil {
		log.Printf("Error getting event: %v", err)
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		err = utils.EncodeErrorResponse(w, ErrEventNotFound)
		if err != nil {
			log.Printf("Encode error wasn't sent to client: %v", err)
		}
		return
	}

	reaction := model.Reaction{
		EventID:   eventId,
		CreatedBy: createdBy,
		LikeValue: 1,
		CreatedAt: time.Now().UTC(),
	}

	err = h.reactionsDomain.SetReaction(ctx, reaction)
	if err != nil {
		log.Printf("SetReaction error: %v", err)
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
		return
	}

	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNoContent)
}

func (h *Handler) SetDislike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("X-Session-Id")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("No cookie: %v", err)
		return
	}

	hexString := cookie.Value
	sid := model.NewSid(hexString)

	createdBy, err := h.eventDomain.GetInternalUserID(ctx, sid)
	if err != nil { // TODO: кастомная ошибка
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusUnauthorized)
		return
	}

	eventId := r.PathValue("event_id")

	_, err = h.eventDomain.GetEventByID(ctx, eventId)
	if err != nil {
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		err = utils.EncodeErrorResponse(w, ErrEventNotFound)
		if err != nil {
			log.Printf("Encode error wasn't sent to client: %v", err)
		}
		return
	}

	reaction := model.Reaction{
		EventID:   eventId,
		CreatedBy: createdBy,
		LikeValue: -1,
		CreatedAt: time.Now().UTC(),
	}

	err = h.reactionsDomain.SetReaction(ctx, reaction)
	if err != nil {
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
		return
	}

	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNoContent)
}
