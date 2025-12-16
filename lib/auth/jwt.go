package auth

import (
	"fmt"
	"maps"
	"time"

	libenv "github.com/fachrunwira/gin-example/lib/env"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims *jwt.MapClaims) (string, error) {
	newClaims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iss": libenv.GetEnv("APP_URL", "http://localhost"),
	}

	if claims != nil {
		maps.Copy(newClaims, *claims)
	}

	generate := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims)
	key := libenv.GetEnv("APP_KEY", "")
	if key == "" {
		return "", fmt.Errorf("cannot generate token APP_KEY must be set first")
	}
	token, err := generate.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func ValidateToken(signedToken string) (*jwt.Token, map[string]any, error) {
	key := libenv.GetEnv("APP_KEY", "")
	if key == "" {
		return nil, map[string]any{}, fmt.Errorf("cannot validate token APP_KEY must be set first")
	}

	parseToken, err := jwt.Parse(signedToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return key, nil
	})

	if err != nil {
		return nil, map[string]any{}, fmt.Errorf("failed to validate token: %w", err)
	}

	if claims, ok := parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {
		newClaims := map[string]any(claims)

		return parseToken, newClaims, nil
	}

	return nil, map[string]any{}, fmt.Errorf("invalid token")
}
