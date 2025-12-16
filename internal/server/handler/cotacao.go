package handler

import (
	"encoding/json"
	"net/http"

	"client-server-api/internal/server/service"
	"client-server-api/pkg/errors"
	"client-server-api/pkg/models"
)

type CotacaoHandler struct {
	service *service.CotacaoService
}

func NewCotacaoHandler(service *service.CotacaoService) *CotacaoHandler {
	return &CotacaoHandler{
		service: service,
	}
}

func (h *CotacaoHandler) GetCotacao(w http.ResponseWriter, r *http.Request) {
	bid, err := h.service.GetBid(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := models.BidResponse{Bid: bid}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *CotacaoHandler) handleError(w http.ResponseWriter, err error) {
	var appErr *errors.AppError
	if !errors.As(err, &appErr) {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
		return
	}

	status := errors.GetHTTPStatus(appErr)

	h.writeJSON(w, status, map[string]string{
		"error": appErr.Message,
	})
}

func (h *CotacaoHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
	}
}



