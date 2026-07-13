package app

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	iam_v1 "iam/internal/api/iam/v1"
	user_v1 "iam/internal/api/user/v1"
	"iam/internal/config"
	"iam/internal/repository"
	sessionRepo "iam/internal/repository/session"
	userRepo "iam/internal/repository/user"
	"iam/internal/service"
	iamService "iam/internal/service/iam"
	userService "iam/internal/service/user"
	"platform/pkg/cache"
	redisPkg "platform/pkg/cache/redis"
	"platform/pkg/closer"
	"platform/pkg/logger"
)

type diContainer struct {
	pool        *pgxpool.Pool
	redisPool   *redigo.Pool
	redisClient cache.RedisClient
	setOperator cache.SetOperator

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	iamService  service.IamService
	userService service.UserService

	userApi *user_v1.UserApi
	iamApi *iam_v1.IamApi
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Pool(ctx context.Context) *pgxpool.Pool {
	if d.pool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to postgres: %v", err))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping postgres: %v", err))
		}

		closer.AddNamed("pgxpool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pool = pool
	}

	return d.pool
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", config.AppConfig().Redis.Address())
			},
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}

		closer.AddNamed("redis pool", func(ctx context.Context) error {
			d.redisPool.Close()
			return nil
		})
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redisPkg.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.Timeout())
	}

	return d.redisClient
}

func (d *diContainer) SetOperator() cache.SetOperator {
	if d.setOperator == nil {
		// redis client implements SetOperator
		d.setOperator = d.RedisClient()
	}

	return d.setOperator
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepo.NewUserRepository(d.Pool(ctx))
	}

	return d.userRepository
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepo.NewSessionRepository(d.RedisClient(), d.SetOperator())
	}

	return d.sessionRepository
}

func (d *diContainer) IamService(ctx context.Context) service.IamService {
	if d.iamService == nil {
		d.iamService = iamService.NewIamService(d.UserRepository(ctx), d.SessionRepository(ctx))
	}

	return d.iamService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewUserService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) UserApi(ctx context.Context) *user_v1.UserApi {
	if d.userApi == nil {
		d.userApi = user_v1.NewUserApi(d.UserService(ctx))
	}

	return d.userApi
}

func (d *diContainer) IamApi(ctx context.Context) *iam_v1.IamApi {
	if d.iamApi == nil {
		d.iamApi = iam_v1.NewIamApi(d.IamService(ctx), d.UserService(ctx))
	}

	return d.iamApi
}