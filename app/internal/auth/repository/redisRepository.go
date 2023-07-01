package repository

import (
	"auth-svc/config"
	"auth-svc/internal/auth"
	"auth-svc/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	models "auth-svc/internal/models/auth"

	pb "gitlab.axarea.ru/main/aiforypay/package/auth-svc-proto"

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

func (r *AuthRedisRepo) GetSession(
	ctx context.Context,
	accessToken string,
	sessionType uint8,
) (models.Session, error) {
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

func (r *AuthRedisRepo) CheckUUIDValid(
	ctx context.Context,
	req *pb.CheckUUIDValidRequest) (res *pb.CheckUUIDValidResponse, err error) {
	ctx, span := otel.Tracer("").Start(ctx, "authRedisRepo.CheckUUIDValid")
	defer span.End()

	var prefix string
	switch req.GetTypeID() {
	case auth.AdminSessionTypeID:
		prefix = adminUUIDPrefix
	case auth.BrokerSessionTypeID:
		prefix = clientUUIDPrefix
	default:
		return &pb.CheckUUIDValidResponse{}, fmt.Errorf("unknown session type id: %d", req.GetTypeID())
	}

	result, err := r.redisClient.GetDel(ctx, prefix+req.GetUUID()).Result()
	if err == redis.Nil {
		return &pb.CheckUUIDValidResponse{},
			errors.Wrapf(err, "AuthRedisRepo.CheckUUIDValid.Get(). Invalid UUID(%s) %s", req.GetUUID(), req.GetTypeID())
	}
	if err != nil {
		return &pb.CheckUUIDValidResponse{},
			errors.Wrap(err, "AuthRedisRepo.CheckUUIDValid.Get(). Unable to get uuid from storage")
	}

	userID, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return &pb.CheckUUIDValidResponse{},
			errors.Wrapf(
				err,
				"AuthRedisRepo.CheckUUIDValid.ParseInt(). Unable to convert userID(%s) from string to int",
				result,
			)
	}

	return &pb.CheckUUIDValidResponse{UserID: int64(userID)}, nil
}

func (r *AuthRedisRepo) SetUUID(ctx context.Context, uuid string, userID int, sessionType uint8) error {
	ctx, span := otel.Tracer("").Start(ctx, "AuthRedisRepo.SetUUID")
	defer span.End()

	var prefix string
	switch sessionType {
	case auth.AdminSessionTypeID:
		prefix = adminUUIDPrefix
	case auth.BrokerSessionTypeID:
		prefix = clientUUIDPrefix
	default:
		return fmt.Errorf("unknown session type id: %d", sessionType)
	}

	res, err := r.redisClient.Set(
		ctx,
		prefix+uuid,
		strconv.Itoa(userID),
		time.Duration(r.cfg.TTLUUID)*time.Second,
	).Result()
	if err != nil {
		r.log.Error(errors.Wrap(err, "AuthRedisRepo.SetUUID"))
		return err
	}

	r.log.Infof("Key(%s):Value(%d) successfully set in Redis. Result:%v", uuid, userID, res)
	return nil
}
