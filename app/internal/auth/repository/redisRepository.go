package repository

import (
	"auth-svc/config"
	"auth-svc/internal/auth"
	"auth-svc/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

	models "auth-svc/internal/models/auth"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

const (
	clientSessPrefix = "CLIENT_SESSION_"
	adminSessPrefix  = "ADMIN_SESSION_"
	clientUUIDPrefix = "WS_CLIENT_"
	adminUUIDPrefix  = "WS_ADMIN_"
)

type AuthRedisRepo struct {
	redisClient *redis.Client
	log         logger.Logger
	cfg         *config.Config
}

func NewAuthRedisRepo(redisClient *redis.Client, log logger.Logger, cfg *config.Config) auth.RedisRepo {
	return &AuthRedisRepo{redisClient: redisClient, log: log, cfg: cfg}
}

func (r *AuthRedisRepo) GetSession(ctx context.Context, accessToken string, sessionType uint8) (models.Session, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userRedisRepo.GetSession")
	defer span.End()

	var prefix string
	switch sessionType {
	case auth.AdminSessionTypeID:
		prefix = adminSessPrefix
	case auth.BrokerSessionTypeID:
		prefix = clientSessPrefix
	default:
		return models.Session{}, fmt.Errorf("unknown session type id: %d", sessionType)
	}

	var key string
	if iter := r.redisClient.Scan(ctx, 0, prefix+"*_"+accessToken, 1).Iterator(); iter.Next(ctx) {
		key = iter.Val()
	} else {
		return models.Session{}, redis.Nil
	}

	sessionString, err := r.redisClient.Get(
		ctx,
		key,
	).Result()
	if err != nil {
		return models.Session{}, errors.Wrapf(
			err,
			"user.repository.GetSession.Get(%s)",
			prefix+accessToken,
		)
	}

	var session models.Session
	err = json.Unmarshal([]byte(sessionString), &session)
	if err != nil {
		return models.Session{}, errors.Wrap(err, "user.repository.GetSession.Unmarshal()")
	}
	return session, nil
}
