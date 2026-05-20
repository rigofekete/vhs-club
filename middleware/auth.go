package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rigofekete/vhs-club-mvc/config"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/internal/auth"
)

const (
	UserIDKey   = "userID"
	UserRoleKey = "userRole"
	RoleUser    = "user"
	RoleAdmin   = "admin"
)

// Middlewares
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, err := extractToken(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		if role != RoleUser && role != RoleAdmin {
			_ = c.Error(apperror.ErrInvalidUser)
			c.Abort()
			return
		}
		c.Set(UserIDKey, userID)
		c.Set(UserRoleKey, role)
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, err := extractToken(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		if role != RoleAdmin {
			_ = c.Error(apperror.ErrInvalidAdmin)
			c.Abort()
			return
		}
		c.Set(UserIDKey, userID)
		c.Set(UserRoleKey, role)
		c.Next()
	}
}

// Helpers
func extractToken(c *gin.Context) (uuid.UUID, string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		_ = c.Error(apperror.ErrInvalidHeader)
		return uuid.Nil, "", apperror.ErrInvalidHeader
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
		_ = c.Error(apperror.ErrInvalidHeader)
		return uuid.Nil, "", apperror.ErrInvalidHeader
	}

	tokenString := splitAuth[1]
	userID, role, err := auth.ValidateJWT(tokenString, config.AppConfig.JWTSecret)
	if err != nil {
		_ = c.Error(apperror.ErrInvalidToken)
		return uuid.Nil, "", apperror.ErrInvalidToken
	}
	return userID, role, nil
}

// Exported helper to extract authenticated user ID from the gin Context's object Keys map, set through User/AdminAuth call.
// This will be called in the Create/ReturnRental handler layer methods, to extract the ID from the gin.Context of the request.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := value.(uuid.UUID)
	return id, ok
}
