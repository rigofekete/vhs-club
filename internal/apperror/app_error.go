package apperror

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppError struct {
	Code    int
	Message string
	Fields  map[string]string `json:"fields,omitempty"`
}

// Sentinel Errors
var (
	// General
	ErrBadRequest = errors.New("bad request")
	// User
	ErrUserNotFound        = errors.New("user not found")
	ErrUserFieldValidation = errors.New("invalid user fields")
	ErrUserExists          = errors.New("user already exists")
	ErrUserInvalidPW       = errors.New("invalid password")
	// Tape
	ErrTapeValidation    = errors.New("invalid tape fields")
	ErrTapeExists        = errors.New("tape already exists")
	ErrTapeNotFound      = errors.New("tape not found")
	ErrTapeUpdateRequest = errors.New("bad update tape request")
	// Rentals
	ErrTapeUnavailable   = errors.New("unavailable tape")
	ErrMaxRentalsPerUser = errors.New("cannot rent more tapes")
	// Auth
	ErrInvalidHeader        = errors.New("invalid header")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidIssuer        = errors.New("invalid issuer")
	ErrInvalidUserID        = errors.New("invalid user id")
	ErrInvalidUser          = errors.New("invalid user")
	ErrInvalidAdmin         = errors.New("invalid admin")
	ErrInvalidSigningMethod = errors.New("invalid JWT signing method")
)

type ValidationError struct {
	Fields map[string]string
}

// ValidationError type implements the Error interface, so we can pass it as an error in the handler layer, c.ShouldBindJSON error check
func (e ValidationError) Error() string {
	return "input validation failed"
}

func WrapValidationError(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	fields := make(map[string]string)

	for _, fieldError := range validationErrors {
		field := fieldError.Field()

		switch fieldError.Tag() {
		case "required":
			fields[field] = "This field is required"
		case "email":
			fields[field] = "Must be a valid email address"
		case "alphanum":
			fields[field] = "Must be formed of only letters and/or numbers"
		case "min":
			fields[field] = "Must be at least " + fieldError.Param() + " characters"
		case "max":
			fields[field] = "Must be at most " + fieldError.Param() + " characters"
		default:
			fields[field] = "Invalid value"
		}
	}
	return ValidationError{Fields: fields}
}

func mapErrorToAppError(err error) *AppError {
	// Check if incoming error is a ValidationError
	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		return &AppError{
			Code:    http.StatusBadRequest,
			Message: "Input validation failed",
			Fields:  validationErr.Fields,
		}
	}

	// Check sentinel errors
	switch {
	case errors.Is(err, ErrBadRequest):
		return &AppError{Code: http.StatusBadRequest, Message: "Bad request"}
	case errors.Is(err, ErrUserNotFound):
		return &AppError{Code: http.StatusNotFound, Message: "User not found"}
	case errors.Is(err, ErrUserExists):
		return &AppError{Code: http.StatusConflict, Message: "User already exists in the DB"}
	case errors.Is(err, ErrUserFieldValidation):
		return &AppError{Code: http.StatusUnprocessableEntity, Message: "Invalid user fields"}
	case errors.Is(err, ErrUserInvalidPW):
		return &AppError{Code: http.StatusUnauthorized, Message: "Invalid password"}
	case errors.Is(err, ErrTapeValidation):
		return &AppError{Code: http.StatusUnprocessableEntity, Message: "Invalid tape fields"}
	case errors.Is(err, ErrTapeExists):
		return &AppError{Code: http.StatusConflict, Message: "Tape already exists in the DB"}
	case errors.Is(err, ErrTapeNotFound):
		return &AppError{Code: http.StatusNotFound, Message: "Tape not found"}
	case errors.Is(err, ErrTapeUpdateRequest):
		return &AppError{Code: http.StatusBadRequest, Message: "Tape update request needs at least 1 non nil value"}
	case errors.Is(err, ErrTapeUnavailable):
		return &AppError{Code: http.StatusUnprocessableEntity, Message: "Sorry, all tapes for this movie are currently rented out."}
	case errors.Is(err, ErrMaxRentalsPerUser):
		return &AppError{Code: http.StatusBadRequest, Message: "Unfortunately, you cannot rent more movies at the moment. Please return one of your currently rented tapes."}
	case errors.Is(err, ErrInvalidHeader):
		return &AppError{Code: http.StatusUnauthorized, Message: "Missing or invalid authorization header"}
	case errors.Is(err, ErrInvalidToken):
		return &AppError{Code: http.StatusUnauthorized, Message: "Invalid or expired token"}
	case errors.Is(err, ErrInvalidSigningMethod):
		return &AppError{Code: http.StatusUnauthorized, Message: "Invalid JWT signing method"}
	case errors.Is(err, ErrInvalidIssuer):
		return &AppError{Code: http.StatusUnauthorized, Message: "Invalid issuer"}
	case errors.Is(err, ErrInvalidUserID):
		return &AppError{Code: http.StatusUnauthorized, Message: "Invalid user ID"}
	case errors.Is(err, ErrInvalidAdmin):
		return &AppError{Code: http.StatusForbidden, Message: "Admin access required"}
	case errors.Is(err, ErrInvalidUser):
		return &AppError{Code: http.StatusForbidden, Message: "User/Admin access required"}
	default:
		return &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	}
}

// Middleware factory function
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}
		appErr := mapErrorToAppError(err.Err)

		response := gin.H{
			"error": appErr.Message,
		}

		if appErr.Fields != nil {
			response["fields"] = appErr.Fields
		}

		c.JSON(appErr.Code, response)
	}
}
