package admin

var (
	// Общие ошибки
	ErrInternalServerError = "Internal server error"
	ErrBadRequest          = "Bad request"
	ErrNotFound            = "Resource not found"
	ErrUnauthorized        = "Unauthorized access"
	ErrForbidden           = "Forbidden access"
	ErrConflict            = "Conflict: resource already exists"
	ErrTooManyRequests     = "Too many requests"
	ErrTimeout             = "Request timeout"

	// Ошибки аутентификации и авторизации
	ErrInvalidCredentials = "Invalid credentials"
	ErrTokenExpired       = "Token has expired"
	ErrTokenInvalid       = "Invalid token"
	ErrAccessDenied       = "Access denied"
	ErrUserNotActive      = "User account is not active"

	// Ошибки валидации
	ErrValidationFailed   = "Validation failed"
	ErrInvalidEmail       = "Invalid email address"
	ErrInvalidPassword    = "Invalid password"
	ErrInvalidUsername    = "Invalid username"
	ErrFieldRequired      = "Field is required"
	ErrFieldTooShort      = "Field is too short"
	ErrFieldTooLong       = "Field is too long"
	ErrInvalidInputFormat = "Invalid input format"
	ErrVerificationFailed = "Verificatio failed"

	// Ошибки, связанные с пользователями
	ErrUserAlreadyExists = "User already exists"
	ErrUserNotFound      = "User not found"
	ErrUserNotCreated    = "Failed to create user"
	ErrUserNotUpdated    = "Failed to update user"
	ErrUserNotDeleted    = "Failed to delete user"

	// Ошибки, связанные с ролями и разрешениями
	ErrRoleNotFound      = "Role not found"
	ErrPermissionDenied  = "Permission denied"
	ErrInvalidRole       = "Invalid role"
	ErrRoleAlreadyExists = "Role already exists"

	// Ошибки, связанные с файлами и загрузкой
	ErrFileTooLarge     = "File is too large"
	ErrInvalidFileType  = "Invalid file type"
	ErrFileUploadFailed = "File upload failed"
	ErrFileNotFound     = "File not found"

	// Ошибки, связанные с базой данных
	ErrDatabaseConnection = "Database connection error"
	ErrDatabaseQuery      = "Database query error"
	ErrDatabaseInsert     = "Database insert error"
	ErrDatabaseUpdate     = "Database update error"
	ErrDatabaseDelete     = "Database delete error"

	// Ошибки, связанные с внешними сервисами
	ErrExternalServiceUnavailable = "External service unavailable"
	ErrExternalServiceTimeout     = "External service timeout"
	ErrExternalServiceError       = "External service error"

	// Ошибки, связанные с конфигурацией
	ErrConfigNotFound   = "Configuration not found"
	ErrInvalidConfig    = "Invalid configuration"
	ErrConfigLoadFailed = "Failed to load configuration"

	// Ошибки, связанные с сессиями
	ErrSessionExpired  = "Session has expired"
	ErrSessionInvalid  = "Invalid session"
	ErrSessionNotFound = "Session not found"

	// Ошибки, связанные с API
	ErrAPINotFound          = "API endpoint not found"
	ErrAPIMethodNotAllowed  = "Method not allowed"
	ErrAPIRateLimitExceeded = "API rate limit exceeded"
)
