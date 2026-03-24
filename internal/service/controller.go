package service

import "github.com/Murolando/ftem_wallet/pkg/entities"

// Controller основной интерфейс контроллера
type Controller interface {
	AuthWallet(mnemonicWords [12]string, password string) entities.ResultString
	GenerateWallet(password string) entities.ResultString
	SendETH() entities.ResultString
	ShowBalance() entities.ResultString
}
