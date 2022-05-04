package errors

import "aspire/models"

func New(code string, message string, detail string) models.ErrorResponse {
	return models.ErrorResponse{Code: code, Message: message, Detail: detail}
}
