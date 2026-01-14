package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production"
	}
	return []byte(secret)
}

func GenerateToken(userID uuid.UUID, email, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "En-tête d'autorisation requis",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Token Bearer requis",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Token invalide ou expiré",
			})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Données du token invalides",
			})
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		return next(c)
	}
}
