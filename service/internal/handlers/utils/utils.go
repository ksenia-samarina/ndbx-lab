package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

func WriteSessionResponse(w http.ResponseWriter, hex string, ttl time.Duration, status int) {
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "X-Session-Id",
		Value:    hex,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(ttl.Seconds()),
	})
	w.WriteHeader(status)
}

func EncodeErrorResponse(w http.ResponseWriter, err error) {
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": err.Error(),
	})
}
