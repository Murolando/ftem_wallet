package service

import "github.com/Murolando/ftem_wallet/pkg/entities"

// Controller основной интерфейс контроллера
type Controller interface {
	AuthWallet() entities.ResultString
	GenerateWallet() entities.ResultString
	SendETH() entities.ResultString
	ShowHistory() entities.ResultString
	ShowBalance() entities.ResultString
}
