package auth

import (
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "vhsclub-access"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, role string, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)

	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    string(TokenTypeAccess),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, string, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperror.ErrInvalidSigningMethod
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil || !token.Valid {
		// NOTE: maybe return a sentinel error here?
		return uuid.Nil, "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, "", err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, "", apperror.ErrInvalidIssuer
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, "", apperror.ErrInvalidUserID
	}
	return id, claims.Role, nil
}
