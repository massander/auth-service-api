package core

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	Token     string    `db:"token"`
	UserID    uuid.UUID `db:"user_id"`
	Rovoked   bool      `db:"rovoked"`
	ClientIP  string    `db:"client_ip"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AccessToken struct {
	ID        uuid.UUID `db:"id"`
	ParentID  uuid.UUID `db:"parent_id"`
	UserID    uuid.UUID `db:"user_id"`
	Revoked   bool      `db:"revoked"`
	ClientIP  string    `db:"client_ip"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TokenPair struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}
