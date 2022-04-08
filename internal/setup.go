package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"go.uber.org/multierr"

	"github.com/Harzu/exchange-rate-test-task/internal/config"
	"github.com/Harzu/exchange-rate-test-task/internal/controllers"
	"github.com/Harzu/exchange-rate-test-task/internal/jobs"
	"github.com/Harzu/exchange-rate-test-task/internal/repositories"
	"github.com/Harzu/exchange-rate-test-task/internal/services"
	"github.com/Harzu/exchange-rate-test-task/internal/system/connections/pg"
	redisclient "github.com/Harzu/exchange-rate-test-task/internal/system/connections/redis"
	"github.com/Harzu/exchange-rate-test-task/internal/system/logger"
)

type app struct {
	ctx           context.Context
	config        *config.Config
	logger        *zerolog.Logger
	pgPool        *pgxpool.Pool
	redisClient   *redis.Client
	jobsScheduler *cron.Cron
	httpServer    *http.Server
}

func NewApp(ctx context.Context) (*app, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	appLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	pgPool, err := pg.NewPool(ctx, cfg.PGDSN)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisclient.New(cfg.Redis, ctx)
	if err != nil {
		return nil, err
	}

	repoContainer := repositories.New(pgPool)
	servicesContainer := services.NewContainer(cfg, appLogger, redisClient, repoContainer)
	jobsScheduler, err := jobs.NewScheduler(cfg.Jobs, appLogger, servicesContainer)
	if err != nil {
		return nil, err
	}
	httpControllers := controllers.NewHTTPContainer(appLogger, servicesContainer.Rates)

	return &app{
		ctx:           ctx,
		logger:        appLogger,
		pgPool:        pgPool,
		config:        cfg,
		jobsScheduler: jobsScheduler,
		redisClient:   redisClient,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: httpControllers.Mux(),
		},
	}, nil
}

func (a *app) Start() {
	a.jobsScheduler.Start()
	a.logger.Info().Str("http_port", a.config.Port).Msg("start http server")
	if err := a.httpServer.ListenAndServe(); err != nil {
		a.logger.Error().Err(err).Msg("failed to start http server")
	}
}

func (a *app) Shutdown() error {
	a.jobsScheduler.Stop()
	err := multierr.Combine(
		a.httpServer.Shutdown(a.ctx),
		a.redisClient.Close(),
	)
	a.pgPool.Close()
	return err
}
