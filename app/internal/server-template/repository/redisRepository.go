package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"server-template/config"
	server_template "server-template/internal/server-template"
	"server-template/pkg/logger"

	"server-template/internal/models/auth"
	"server-template/internal/models/templates"

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

type ServerTeamplateRedisRepo struct {
	redisClient *redis.Client
	log         logger.Logger
	cfg         *config.Config
}

func NewServerTeamplateRedisRepo(redisClient *redis.Client, log logger.Logger, cfg *config.Config) server_template.RedisRepo {
	return &ServerTeamplateRedisRepo{redisClient: redisClient, log: log, cfg: cfg}
}

func (r *ServerTeamplateRedisRepo) GetSession(ctx context.Context, accessToken string, sessionType uint8) (auth.Session, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userRedisRepo.GetSession")
	defer span.End()

	var prefix string
	switch sessionType {
	case server_template.AdminSessionTypeID:
		prefix = adminSessPrefix
	case server_template.BrokerSessionTypeID:
		prefix = clientSessPrefix
	default:
		return auth.Session{}, fmt.Errorf("unknown session type id: %d", sessionType)
	}

	var key string
	if iter := r.redisClient.Scan(ctx, 0, prefix+"*_"+accessToken, 1).Iterator(); iter.Next(ctx) {
		key = iter.Val()
	} else {
		return auth.Session{}, redis.Nil
	}

	sessionString, err := r.redisClient.Get(
		ctx,
		key,
	).Result()
	if err != nil {
		return auth.Session{}, errors.Wrapf(
			err,
			"user.repository.GetSession.Get(%s)",
			prefix+accessToken,
		)
	}

	var session auth.Session
	err = json.Unmarshal([]byte(sessionString), &session)
	if err != nil {
		return auth.Session{}, errors.Wrap(err, "user.repository.GetSession.Unmarshal()")
	}
	return session, nil
}

func (r *ServerTeamplateRedisRepo) HTTPPing(ctx context.Context) (result templates.Templates, err error) {
	result.Pong = "Pong"
	return
}
