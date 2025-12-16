package handler

import (
	"encoding/json"
	"net/http"

	"client-server-api/internal/server/service"
	"client-server-api/pkg/errors"
	"client-server-api/pkg/models"
)

// CotacaoHandler lida com requisições HTTP relacionadas a cotações
type CotacaoHandler struct {
	service *service.CotacaoService
}

// NewCotacaoHandler cria uma nova instância do handler de cotação
func NewCotacaoHandler(service *service.CotacaoService) *CotacaoHandler {
	return &CotacaoHandler{
		service: service,
	}
}

// GetCotacao busca a cotação do dólar e retorna como JSON
func (h *CotacaoHandler) GetCotacao(w http.ResponseWriter, r *http.Request) {
	// Chamar service
	bid, err := h.service.GetBid(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Criar resposta
	response := models.BidResponse{Bid: bid}

	// Retornar JSON
	h.writeJSON(w, http.StatusOK, response)
}

// handleError trata erros e retorna resposta HTTP apropriada
func (h *CotacaoHandler) handleError(w http.ResponseWriter, err error) {
	var appErr *errors.AppError
	if !errors.As(err, &appErr) {
		// Erro genérico não tratado
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// Obter status HTTP apropriado
	status := errors.GetHTTPStatus(appErr)

	// Retornar erro JSON
	h.writeJSON(w, status, map[string]string{
		"error": appErr.Message,
	})
}

// writeJSON escreve uma resposta JSON com status code apropriado
func (h *CotacaoHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Erro ao escrever JSON (já escreveu headers, não pode retornar erro ao cliente)
		// Em produção, deveria logar este erro
	}
}


