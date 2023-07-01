package auth

import (
	models "auth-svc/internal/models/auth"
	"context"
)

type UseCase interface {
	GetSession(ctx context.Context, accessToken string, sessionType uint8) (models.Session, error)
}
