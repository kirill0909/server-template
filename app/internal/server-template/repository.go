package server_template

import (
	"context"
	models "server-template/internal/models/auth"
)

type RedisRepo interface {
	GetSession(ctx context.Context, accessToken string, sessionType uint8) (models.Session, error)
}
