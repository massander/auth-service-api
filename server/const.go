package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const SIGNING_ALGORITHM = "SHA512"

const ACCESS_TOKEN_TTL = time.Minute * 5
const REFRESH_TOKEN_TTL = time.Hour * 24

const JWT_SIGNING_SECRET = "secret_to_change"
const JWT_ISSUER = "auth-service-api"

var JWT_SIGNING_METHOD = jwt.SigningMethodHS512

const MIGRATIONS_FOLDER = "/storage/postgres/migrations"
