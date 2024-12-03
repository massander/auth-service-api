package server

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateAccessToken(now time.Time, id uuid.UUID, userID uuid.UUID, clientIP string) (string, error) {
	const op = "server.tokens.generateAccessToken"

	iat := jwt.NewNumericDate(now)
	exp := jwt.NewNumericDate(now.Add(ACCESS_TOKEN_TTL))

	claims := jwt.MapClaims{
		"iss":       JWT_ISSUER,
		"sub":       userID.String(),
		"iat":       iat,
		"exp":       exp,
		"client_ip": clientIP,
		"jti":       id.String(),
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)

	signedToken, err := token.SignedString([]byte(JWT_SIGNING_SECRET))
	if err != nil {
		return "", fmt.Errorf("%s: filed to sign access token: %w", op, err)
	}

	return signedToken, nil

}
func generateRefreshToken(now time.Time, id uuid.UUID, userID uuid.UUID, clientIP string) (string, error) {
	const op = "server.tokens.generateRefreshToken"

	iat := jwt.NewNumericDate(now)
	exp := jwt.NewNumericDate(now.Add(REFRESH_TOKEN_TTL))

	claims := jwt.MapClaims{
		"iss":       JWT_ISSUER,
		"sub":       userID.String(),
		"iat":       iat,
		"exp":       exp,
		"client_ip": clientIP,
		"jti":       id.String(),
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)

	signedToken, err := token.SignedString([]byte(JWT_SIGNING_SECRET))
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign refresh token: %w", op, err)
	}

	return signedToken, nil
}

func parseToken(token string) (*jwt.Token, error) {
	const op = "server.tokens.parseToken"
	t, err := jwt.ParseWithClaims(
		token,
		&jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(JWT_SIGNING_SECRET), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}
