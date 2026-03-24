package cli

import (
	"strings"

	"github.com/Murolando/ftem_wallet/pkg/entities"
)

// AuthWallet авторизует кошелек (используется только при запуске с флагом)
func (c *CLIController) AuthWallet(mnemonicWords [12]string, password string) entities.ResultString {
	result := c.service.AuthWallet(mnemonicWords, password)
	if !strings.Contains(string(result), "Error") {
		c.isAuthorized = true
	}
	return result
}

// GenerateWallet генерирует новый кошелек
func (c *CLIController) GenerateWallet(password string) entities.ResultString {
	result := c.service.GenerateWallet(password)
	if !strings.Contains(string(result), "Error") {
		c.isAuthorized = true
	}
	return result
}
