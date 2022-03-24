package internal

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Harzu/exchange-rate-test-task/internal/system/connections/mongodb"

	"github.com/rs/zerolog"

	"github.com/Harzu/exchange-rate-test-task/internal/system/logger"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Harzu/exchange-rate-test-task/internal/repositories"

	"github.com/Harzu/exchange-rate-test-task/internal/services/rates"

	"github.com/Harzu/exchange-rate-test-task/internal/controllers"

	"github.com/Harzu/exchange-rate-test-task/internal/config"
	"github.com/Harzu/exchange-rate-test-task/internal/system/connections/pg"
)

type app struct {
	ctx         context.Context
	logger      *zerolog.Logger
	pgPool      *pgxpool.Pool
	mongoClient *mongo.Client
	httpServer  *http.Server
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

	mongoClient, err := mongodb.NewClient(ctx, cfg.MongoURI)
	if err != nil {
		return nil, err
	}

	pgPool, err := pg.NewPool(ctx, cfg.PGDSN)
	if err != nil {
		return nil, err
	}

	repoContainer := repositories.New(pgPool)
	ratesService := rates.NewService(repoContainer)
	httpControllers := controllers.NewHTTPContainer(ratesService)

	return &app{
		ctx:         ctx,
		logger:      appLogger,
		pgPool:      pgPool,
		mongoClient: mongoClient,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: httpControllers.Mux(),
		},
	}, nil
}

func (a *app) Start() {
	if err := a.httpServer.ListenAndServe(); err != nil {
		a.logger.Error().Err(err).Msg("failed to start http server")
	}
}

// todo: multierr
func (a *app) Shutdown() error {
	if err := a.httpServer.Shutdown(a.ctx); err != nil {
		return err
	}
	if err := a.mongoClient.Disconnect(a.ctx); err != nil {
		return err
	}
	return nil
}
