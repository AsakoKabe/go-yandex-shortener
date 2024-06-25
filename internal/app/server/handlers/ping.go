package handlers

import (
	"encoding/json"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"net/http"
)

type PingHandler struct {
	pingService service.PingService
}

func NewPingHandler(pingService service.PingService) *PingHandler {
	return &PingHandler{pingService: pingService}
}

func (h *PingHandler) healthDB(w http.ResponseWriter, r *http.Request) {
	err := h.pingService.PingDB(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal(map[string]string{"name": err.Error()})
		http.Error(w, string(message), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}
