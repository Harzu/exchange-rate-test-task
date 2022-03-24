package services

import (
	"github.com/Harzu/exchange-rate-test-task/internal/repositories"
	"github.com/Harzu/exchange-rate-test-task/internal/services/rates"
)

type Container struct {
	Rates *rates.Service
}

func New(repoContainer *repositories.Container) *Container {
	return &Container{
		Rates: rates.NewService(repoContainer.Rates),
	}
}
