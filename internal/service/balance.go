package service

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Murolando/ftem_wallet/pkg/entities"
)

// ShowBalance показывает баланс кошелька
func (s *Service) ShowBalance() entities.ResultString {
	var result strings.Builder
	result.WriteString("Балансы счетов кошелька:\n")

	totalBalance := big.NewInt(0)

	for i, address := range s.currentWallet.Addresses {
		balance, err := s.ethClient.GetBalance(address)
		if err != nil {
			result.WriteString(fmt.Sprintf("Адрес %d (%s): Ошибка получения баланса - %v\n", i, address, err))
			continue
		}

		// Конвертируем wei в ETH (1 ETH = 10^18 wei)
		ethBalance := new(big.Float).SetInt(balance)
		ethBalance.Quo(ethBalance, big.NewFloat(1e18))

		result.WriteString(fmt.Sprintf("Адрес %d (%s): %s ETH\n", i, address, ethBalance.Text('f', 18)))
		totalBalance.Add(totalBalance, balance)
	}

	// Общий баланс
	totalEthBalance := new(big.Float).SetInt(totalBalance)
	totalEthBalance.Quo(totalEthBalance, big.NewFloat(1e18))
	result.WriteString(fmt.Sprintf("\nОбщий баланс: %s ETH", totalEthBalance.Text('f', 18)))

	return entities.ResultString(result.String())
}

// SendETH отправляет ETH с указанного адреса на другой адрес
func (s *Service) SendETH(fromAddress, toAddress, amountETH string) entities.ResultString {
	fromIndex := -1
	for i, addr := range s.currentWallet.Addresses {
		if strings.EqualFold(addr, fromAddress) {
			fromIndex = i
			break
		}
	}

	if fromIndex == -1 {
		return entities.ResultString(fmt.Sprintf("Ошибка: адрес %s не найден в кошельке", fromAddress))
	}

	privateKey := s.currentWallet.GetPrivateKey(fromIndex)
	if privateKey == "" {
		return entities.ResultString("Ошибка: не удалось получить приватный ключ")
	}

	// Конвертировать количество ETH в wei
	amountFloat := new(big.Float)
	amountFloat, ok := amountFloat.SetString(amountETH)
	if !ok {
		return entities.ResultString("Ошибка: неверный формат суммы")
	}

	weiPerEth := new(big.Float).SetInt(big.NewInt(1e18))
	amountWei := new(big.Float).Mul(amountFloat, weiPerEth)

	amountWeiInt := new(big.Int)
	amountWei.Int(amountWeiInt)

	// Получить цену газа
	gasPrice, err := s.ethClient.GetGasPrice()
	if err != nil {
		return entities.ResultString(fmt.Sprintf("Ошибка получения цены газа: %v", err))
	}

	// Установить лимит газа (стандартный для простого перевода ETH)
	gasLimit := uint64(21000)

	// Проверить баланс перед отправкой транзакции
	balance, err := s.ethClient.GetBalance(fromAddress)
	if err != nil {
		return entities.ResultString(fmt.Sprintf("Ошибка получения баланса: %v", err))
	}

	// Рассчитать общую стоимость транзакции (сумма + газ)
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	totalCost := new(big.Int).Add(amountWeiInt, gasCost)

	// Проверить, достаточно ли средств
	if balance.Cmp(totalCost) < 0 {
		// Конвертируем в ETH для удобного отображения
		balanceETH := new(big.Float).SetInt(balance)
		balanceETH.Quo(balanceETH, big.NewFloat(1e18))

		totalCostETH := new(big.Float).SetInt(totalCost)
		totalCostETH.Quo(totalCostETH, big.NewFloat(1e18))

		gasCostETH := new(big.Float).SetInt(gasCost)
		gasCostETH.Quo(gasCostETH, big.NewFloat(1e18))

		return entities.ResultString(fmt.Sprintf(
			"Ошибка: недостаточно средств для транзакции\n"+
				"Баланс: %s ETH\n"+
				"Требуется: %s ETH (сумма: %s ETH + газ: %s ETH)",
			balanceETH.Text('f', 18),
			totalCostETH.Text('f', 18),
			amountETH,
			gasCostETH.Text('f', 18),
		))
	}
	// Получить nonce
	nonce, err := s.ethClient.GetTransactionCount(fromAddress)
	if err != nil {
		return entities.ResultString(fmt.Sprintf("Ошибка получения nonce: %v", err))
	}

	// Получить chainID
	chainID, err := s.ethClient.ChainID()
	if err != nil {
		return entities.ResultString(fmt.Sprintf("Ошибка получения chainID: %v", err))
	}

	// Подписать транзакцию
	signedTx, err := s.ethClient.SignTransaction(
		privateKey,
		toAddress,
		amountWeiInt,
		gasLimit,
		gasPrice,
		nonce,
		nil, // data пустая для простого перевода ETH
		chainID,
	)
	if err != nil {
		return entities.ResultString(fmt.Sprintf("Ошибка подписания транзакции: %v", err))
	}

	// Отправить транзакцию
	txHash, err := s.ethClient.SendTransaction(signedTx)
	if err != nil {
		return entities.ResultString(fmt.Sprintf("Ошибка отправки транзакции: %v", err))
	}

	return entities.ResultString(fmt.Sprintf(
		"Транзакция успешно отправлена!\n"+
			"От: %s\n"+
			"Кому: %s\n"+
			"Сумма: %s ETH\n"+
			"Хеш транзакции: %s\n"+
			"Nonce: %d\n"+
			"Gas Price: %s wei",
		fromAddress, toAddress, amountETH, txHash, nonce, gasPrice.String(),
	))
}
