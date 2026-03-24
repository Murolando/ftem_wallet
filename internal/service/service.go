package service

import (
	"github.com/Murolando/ftem_wallet/pkg/config"
	"github.com/Murolando/ftem_wallet/pkg/entities"
	"github.com/Murolando/ftem_wallet/pkg/lib/clients/eth"
)

// Service основная структура сервиса
type Service struct {
	serviceConfig   *config.ServiceConfig
	ethClient       eth.Etherium
	currentWallet   *entities.Wallet
}

func New(
	cfg *config.ServiceConfig,
) *Service {
	// Инициализируем Ethereum клиент с URL из конфигурации
	ethClient := eth.NewEtheriumClient(cfg.Host)
	
	return &Service{
		serviceConfig: cfg,
		ethClient:     ethClient,
	}
}

// SetCurrentWallet устанавливает текущий кошелек
func (s *Service) SetCurrentWallet(wallet *entities.Wallet) {
	s.currentWallet = wallet
}

// GetCurrentWallet возвращает текущий кошелек
func (s *Service) GetCurrentWallet() *entities.Wallet {
	return s.currentWallet
}
