package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"samarina/ndbx/internal/domains/validator"
	"samarina/ndbx/internal/handlers/utils"
	"samarina/ndbx/internal/model"
	"strconv"
	"time"
)

type Handler struct {
	userDomain      userDomain
	validatorDomain validatorDomain
	ttl             time.Duration
}

func New(userDomain userDomain, validatorDomain validatorDomain, ttl time.Duration) *Handler {
	return &Handler{
		userDomain:      userDomain,
		validatorDomain: validatorDomain,
		ttl:             ttl,
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
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "limit"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}
		if query.Get("limit") == "" {
			limit = 10
		}

		offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if query.Get("offset") != "" && (err != nil || offset < 0) {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "offset"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		name := query.Get("name")
		id := query.Get("id")

		users, err := h.userDomain.GetUsers(ctx, id, name, uint64(limit), uint64(offset))

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
			err := h.userDomain.UpdateUserSession(ctx, sid, h.ttl)
			if err != nil {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
				return
			}
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "username"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}
		if user.FullName == "" {
			err := h.userDomain.UpdateUserSession(ctx, sid, h.ttl)
			if err != nil {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
				return
			}
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "full_name"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}
		if user.Password == "" {
			err := h.userDomain.UpdateUserSession(ctx, sid, h.ttl)
			if err != nil {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
				return
			}
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
			err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: "password"})
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		_, err := h.userDomain.GetUserByUsername(ctx, user.Username)
		if err == nil { // пользователь уже существует
			err = h.userDomain.UpdateUserSession(ctx, sid, h.ttl)
			if err != nil {
				utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
				return
			}
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusConflict)
			err = utils.EncodeErrorResponse(w, ErrUserAlreadyExists)
			if err != nil {
				log.Printf("Encode error wasn't sent to client: %v", err)
			}
			return
		}

		newSid, err := h.userDomain.RegisterUser(ctx, user, h.ttl)
		if err != nil {
			utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
			return
		}
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
	user, err := h.userDomain.GetUserByUserID(ctx, id)
	if err != nil { // TODO: кастомная ошибка что нет пользака
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		err = utils.EncodeErrorResponse(w, ErrUserNotFound)
		if err != nil {
			log.Printf("Encode error wasn't sent to client: %v", err)
		}
		return
	}
	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
	err = h.userDomain.UpdateUserSession(ctx, sid, h.ttl)
	if err != nil {
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
		return
	}
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
	_, err := h.userDomain.GetUserByUserID(ctx, id)
	if err != nil { // TODO: кастомная ошибка что нет пользака
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusNotFound)
		err = utils.EncodeErrorResponse(w, ErrUserNotFound)
		if err != nil {
			log.Printf("Encode error wasn't sent to client: %v", err)
		}
		return
	}

	query := r.URL.Query()

	var target *validator.ErrInvalidFieldName
	filter, err := h.validatorDomain.ValidateParams(query)
	if errors.As(err, &target) {
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusBadRequest)
		err = utils.EncodeErrorResponse(w, &ErrInvalidFieldName{Field: target.Field})
		if err != nil {
			log.Printf("Encode error wasn't sent to client: %v", err)
		}
	}
	filter.User = id

	events, err := h.userDomain.GetEvents(ctx, filter)
	if err != nil {
		utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"events": events,
		"count":  len(events),
	}
	utils.WriteSessionResponse(w, sid.HexString, h.ttl, http.StatusOK)
	_ = json.NewEncoder(w).Encode(&resp)
}
