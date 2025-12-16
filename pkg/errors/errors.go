package errors

import (
	"context"
	stderrors "errors"
	"net/http"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func ErroTimeout(message string) *AppError {
	return &AppError{
		Code:    "TIMEOUT",
		Message: message,
		Err:     nil,
	}
}

// ErroTimeoutContext cria erro de timeout a partir de um contexto
func ErroTimeoutContext(operation string, err error) *AppError {
	if err == context.DeadlineExceeded {
		return &AppError{
			Code:    "TIMEOUT",
			Message: "Timeout ao executar operação: " + operation,
			Err:     err,
		}
	}
	return &AppError{
		Code:    "TIMEOUT",
		Message: "Timeout ao executar operação: " + operation,
		Err:     err,
	}
}

// ErroAPI cria erro de chamada à API externa
func ErroAPI(err error) *AppError {
	return &AppError{
		Code:    "API_ERROR",
		Message: "Erro ao chamar API externa",
		Err:     err,
	}
}

// ErroDatabase cria erro de banco de dados
func ErroDatabase(err error) *AppError {
	return &AppError{
		Code:    "DATABASE_ERROR",
		Message: "Erro ao acessar banco de dados",
		Err:     err,
	}
}

// ErroValidacao cria erro de validação
func ErroValidacao(message string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Err:     nil,
	}
}

// ErroInterno cria erro interno genérico
func ErroInterno(err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Erro interno do servidor",
		Err:     err,
	}
}

// ErroNotFound cria erro de recurso não encontrado
func ErroNotFound(resource string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: resource + " não encontrado",
		Err:     nil,
	}
}

// GetHTTPStatus mapeia o código do erro para status HTTP apropriado
func GetHTTPStatus(err error) int {
	var appErr *AppError
	if !As(err, &appErr) {
		return http.StatusInternalServerError
	}

	switch appErr.Code {
	case "TIMEOUT":
		return http.StatusGatewayTimeout
	case "API_ERROR":
		return http.StatusBadGateway
	case "DATABASE_ERROR":
		return http.StatusInternalServerError
	case "VALIDATION_ERROR":
		return http.StatusBadRequest
	case "NOT_FOUND":
		return http.StatusNotFound
	case "INTERNAL_ERROR":
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// As verifica se o erro é do tipo AppError (usa errors.As do pacote padrão)
func As(err error, target **AppError) bool {
	return stderrors.As(err, target)
}
