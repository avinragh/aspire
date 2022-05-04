package errors

const (
	ErrInputValidationCode = "INPUT_VALIDATION_ERR"
	ErrInternalServerCode  = "INTERNAL_SERVER_ERR"
	ErrUnauthorizedCode    = "UNAUTHORIZED_ERR"
	ErrForbiddenCode       = "FORBIDDEN_ERR"
	ErrConflictCode        = "RESOURCE_CONFLICT_ERR"
)

const (
	ErrInputValidationMessage = "Request failed Validation"
	ErrInternalServerMessage  = "Internal Server Error - please contact support"
	ErrAuthenticationMessage  = "User is not authenticated"
	ErrForbiddenMessage       = "User cannot perform the action"
	ErrConflictMessage        = "The resource Already exists"
)
