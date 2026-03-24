package cli

import "github.com/Murolando/ftem_wallet/pkg/entities"

// ShowBalance показывает баланс кошелька
func (c *CLIController) ShowBalance() entities.ResultString {
	return c.service.ShowBalance()
}

// SendETH отправляет ETH
func (c *CLIController) SendETH(fromAddress, toAddress, amount string) entities.ResultString {
	return c.service.SendETH(fromAddress, toAddress, amount)
}
