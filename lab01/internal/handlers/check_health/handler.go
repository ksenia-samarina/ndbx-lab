package check_health

import (
	"encoding/json"
	"log"
	"net/http"
	"samarina/ndbx/internal/dto/response"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	healthCheck := response.HealthCheck{
		Status: response.OkHealthCheck,
	}

	err := json.NewEncoder(w).Encode(healthCheck)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
