package jobs

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"

	"github.com/Harzu/exchange-rate-test-task/internal/config"
	"github.com/Harzu/exchange-rate-test-task/internal/services"
)

func NewScheduler(
	cfg *config.Jobs,
	logger *zerolog.Logger,
	services *services.Container,
) (*cron.Cron, error) {
	scheduler := cron.New(cron.WithLogger(newCronLogger(logger)))

	uploadRatesJob := newUploadRatesJob(cfg, logger, services.Locker, services.Rates)
	if _, err := scheduler.AddJob(uploadRatesJob.Spec(), uploadRatesJob); err != nil {
		return nil, fmt.Errorf("failed to add simple job: %w", err)
	}

	return scheduler, nil
}
