package service

import (
	"github.com/Murolando/ftem_wallet/internal/repository"
	"github.com/Murolando/ftem_wallet/pkg/config"
)

// Service основная структура сервиса
type Service struct {
	serviceConfig *config.ServiceConfig
	repository    repository.Repository
}

func New(
	r repository.Repository,
	cfg *config.ServiceConfig,

) *Service {
	return &Service{
		repository:    r,
		serviceConfig: cfg,
	}
}
