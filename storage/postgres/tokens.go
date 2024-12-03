package postgres

import (
	"auth-service-api/core"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) FindTokens(ctx context.Context, refreshJTI, userID uuid.UUID) (core.TokenPair, error) {
	const op = "storage.postgres.tokens.FindTokens"

	query := "select * from refresh_tokens where id=$1 and user_id=$2"
	var accessToken core.AccessToken
	if err := s.pool.QueryRow(ctx, query, refreshJTI, userID).Scan(&accessToken); err != nil {
		return core.TokenPair{}, fmt.Errorf("%s: %w", op, err)
	}

	query = "select * from access_token where parent_id=$1 and user_id=$2"
	var refreshToken core.RefreshToken
	if err := s.pool.QueryRow(ctx, query, accessToken.ParentID, userID).Scan(&refreshToken); err != nil {
		return core.TokenPair{}, fmt.Errorf("%s: %w", op, err)
	}

	return core.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Storage) RevokeTokens(ctx context.Context, accessTokenID, refreshTokenID uuid.UUID) error {
	const op = "storage.postgres.tokens.RevokeTokens"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	query := `update access_tokens set revoked=true where id=$1`
	_, err = tx.Exec(ctx, query, accessTokenID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	query = `update refresh_tokens set revoked=true where id=$1`
	_, err = tx.Exec(ctx, query, refreshTokenID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveTokens(ctx context.Context, tokens core.TokenPair) error {
	const op = "storage.postgres.tokens.SaveTokens"

	userID := tokens.AccessToken.UserID
	clientIP := tokens.AccessToken.ClientIP

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	query := `insert into refresh_tokens(id, "token", user_id, client_ip) values($1,$2,$3,$4)`
	_, err = tx.Exec(ctx, query,
		tokens.RefreshToken.ID,
		tokens.RefreshToken.Token,
		userID,
		clientIP,
	)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	query = `insert into access_tokens(id, parent_id, user_id, client_ip) values($1,$2,$3,$4)`
	_, err = tx.Exec(ctx, query,
		tokens.AccessToken.ID,
		tokens.RefreshToken.ID,
		userID,
		clientIP,
	)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}
