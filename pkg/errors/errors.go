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

func ErroAPI(err error) *AppError {
	return &AppError{
		Code:    "API_ERROR",
		Message: "Erro ao chamar API externa",
		Err:     err,
	}
}

func ErroDatabase(err error) *AppError {
	return &AppError{
		Code:    "DATABASE_ERROR",
		Message: "Erro ao acessar banco de dados",
		Err:     err,
	}
}

func ErroValidacao(message string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Err:     nil,
	}
}

func ErroInterno(err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Erro interno do servidor",
		Err:     err,
	}
}

func ErroNotFound(resource string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: resource + " não encontrado",
		Err:     nil,
	}
}

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

func As(err error, target **AppError) bool {
	return stderrors.As(err, target)
}
