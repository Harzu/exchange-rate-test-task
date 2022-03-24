package rates

import (
	"context"

	"github.com/Harzu/exchange-rate-test-task/internal/repositories"

	"github.com/Harzu/exchange-rate-test-task/internal/repositories/rates"

	"github.com/Harzu/exchange-rate-test-task/internal/entities"
)

type Service struct {
	repository *rates.Repository
}

func NewService(repoContainer *repositories.Container) *Service {
	return &Service{
		repository: repoContainer.Rates,
	}
}

func (s *Service) GetRate(ctx context.Context, sourceSym, targetSym []string) ([]entities.Pair, error) {
	return []entities.Pair{}, nil
}
