package server

import (
	"auth-service-api/core"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Server) handleGetToken(c *fiber.Ctx) error {
	userID := c.Query("user_id", "")
	log.Debug("UserID:", userID)

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETR", "user_id parameter must be valid UUID string")
	}

	now := time.Now()
	accessJTI := uuid.New()
	refreshJTI := uuid.New()
	clientIP := c.IP()

	accessToken, err := generateAccessToken(now, accessJTI, parsedUserID, clientIP)
	if err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	}

	refreshToken, err := generateRefreshToken(now, refreshJTI, parsedUserID, clientIP)
	if err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	}

	// TODO: bcrypt(sha256(token , salt)) to avoid go out of 72 length
	// hashbytes, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.MinCost)
	// if err != nil {
	// 	log.Errorf("generate hash: %w", err)
	// 	return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	// }

	tokenPair := core.TokenPair{
		AccessToken: core.AccessToken{
			ID:       accessJTI,
			ParentID: refreshJTI,
			UserID:   parsedUserID,
			ClientIP: clientIP,
		},
		RefreshToken: core.RefreshToken{
			ID:       refreshJTI,
			Token:    refreshToken,
			UserID:   parsedUserID,
			ClientIP: clientIP,
		},
	}

	if err := s.storage.SaveTokens(c.Context(), tokenPair); err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	}

	return okResponse(c, Tokens{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (s *Server) handleRefresh(c *fiber.Ctx) error {
	var payload Tokens
	if err := c.BodyParser(&payload); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", "failed to parse request body")
	}

	refreshToken, err := parseToken(payload.RefreshToken)
	if err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", "failed to parse refresh token")
	}

	accessToken, err := parseToken(payload.AccessToken)
	if err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", "failed to parse access token")
	}

	// Refresh Token Claims
	refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		log.Debug(refreshClaims)
		return errorResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", nil)
	}

	// Access Token Claims
	accessClaims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Debug(accessClaims)
		return errorResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", nil)
	}

	clientIP := c.IP()
	if refreshClaims["client_ip"] != clientIP {
		// TODO:send email notification
		log.Warn("Client IP has been chaned")
	}

	accesssJTI, err := uuid.Parse(accessClaims["jti"].(string))
	if err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	}

	refreshJTI, err := uuid.Parse(refreshClaims["jti"].(string))
	if err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	}

	if err := s.storage.RevokeTokens(c.Context(), accesssJTI, refreshJTI); err != nil {
		log.Error(err)
		return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	}

	return s.handleGetToken(c)
	// if err != nil {
	// 	return errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", nil)
	// }

	// return nil
}
