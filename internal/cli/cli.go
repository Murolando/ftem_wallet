package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Murolando/ftem_wallet/internal/service"
	"github.com/Murolando/ftem_wallet/pkg/config"
	"github.com/Murolando/ftem_wallet/pkg/entities"
)

type CLIController struct {
	service      *service.Service
	config       *config.Config
	isAuthorized bool
	scanner      *bufio.Scanner
}

func NewCLIController(cfg *config.Config) *CLIController {
	svc := service.New(&cfg.ServiceConfig)
	return &CLIController{
		service: svc,
		config:  cfg,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (c *CLIController) Init(ctx context.Context) error {
	fmt.Println("🔐 Ethereum Wallet CLI")
	fmt.Println("======================")

	// Проверяем, есть ли мнемоника в флагах
	if c.config.CLIConfig.Mnemonic != "" {
		fmt.Println("📝 Обнаружена мнемоническая фраза, выполняется авторизация...")
		result := c.authWalletFromConfig()
		fmt.Println(string(result))

		if strings.Contains(string(result), "Error") {
			fmt.Println("❌ Ошибка авторизации!")
			return fmt.Errorf("failed to authorize wallet")
		}

		c.isAuthorized = true
		fmt.Println("✅ Авторизация успешна!")
	} else {
		fmt.Println("ℹ️  Для начала работы необходимо создать или импортировать кошелек")
	}

	return nil
}

func (c *CLIController) Run(ctx context.Context) error {

	var (
		menu    = c.showAuthorizedMenu
		choiceF = c.handleAuthorizedChoice
	)
	if !c.isAuthorized {
		menu = c.showUnauthorizedMenu
		choiceF = c.handleUnauthorizedChoice
	} 

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		menu()
		choice := c.readInput("Выберите опцию: ")
		if !choiceF(choice) {
			break
		}

		fmt.Println()
	}

	fmt.Println("👋 До свидания!")
	return nil
}

func (c *CLIController) showUnauthorizedMenu() {
	fmt.Println("\n📋 Доступные действия:")
	fmt.Println("1. 🆕 Создать новый кошелек")
	fmt.Println("2. ❌ Выход")
}

func (c *CLIController) showAuthorizedMenu() {
	fmt.Println("\n📋 Доступные действия:")
	fmt.Println("1. 👁️  Показать адреса кошелька")
	fmt.Println("2. 💰 Показать баланс")
	fmt.Println("3. 📤 Отправить ETH")
	fmt.Println("4. ❌ Выход")
}

func (c *CLIController) handleUnauthorizedChoice(choice string) bool {
	switch choice {
	case "1":
		c.generateWallet()
	case "2":
		return false
	default:
		fmt.Println("❌ Неверный выбор, попробуйте снова")
	}
	return true
}

func (c *CLIController) handleAuthorizedChoice(choice string) bool {
	switch choice {
	case "1":
		c.showWalletAddresses()
	case "2":
		c.showBalance()
	case "3":
		c.sendETH()
	case "4":
		return false
	default:
		fmt.Println("❌ Неверный выбор, попробуйте снова")
	}
	return true
}

func (c *CLIController) readInput(prompt string) string {
	fmt.Print(prompt)
	c.scanner.Scan()
	return strings.TrimSpace(c.scanner.Text())
}

func (c *CLIController) authWalletFromConfig() entities.ResultString {
	// Парсим мнемоническую фразу
	words := strings.Fields(c.config.CLIConfig.Mnemonic)
	if len(words) != 12 {
		return entities.ResultString("Error: Мнемоническая фраза должна содержать 12 слов")
	}

	var mnemonicArray [12]string
	copy(mnemonicArray[:], words)

	return c.service.AuthWallet(mnemonicArray, c.config.CLIConfig.Password)
}

func (c *CLIController) generateWallet() {
	fmt.Println("\n🆕 Создание нового кошелька")
	password := c.readInput("Введите пароль для кошелька (или оставьте пустым): ")

	result := c.service.GenerateWallet(password)
	fmt.Println(string(result))

	if !strings.Contains(string(result), "Error") {
		c.isAuthorized = true
		fmt.Println("✅ Кошелек успешно создан!")
		fmt.Println("⚠️  ВАЖНО: Сохраните мнемоническую фразу в безопасном месте!")
	}
}

func (c *CLIController) showWalletAddresses() {
	fmt.Println("\n👁️  Адреса кошелька:")
	wallet := c.service.GetCurrentWallet()

	for i, addr := range wallet.Addresses {
		fmt.Printf("Кошелек %d: %s\n", i, addr)
	}
}

func (c *CLIController) showBalance() {
	fmt.Println("\n💰 Получение баланса...")
	result := c.service.ShowBalance()
	fmt.Println(string(result))
}

func (c *CLIController) sendETH() {
	fmt.Println("\n📤 Отправка ETH")

	wallet := c.service.GetCurrentWallet()

	// Показываем доступные адреса
	fmt.Println("Доступные адреса для отправки:")
	for i, addr := range wallet.Addresses {
		fmt.Printf("%d. %s\n", i, addr)
	}

	// Цикл для выбора адреса отправителя
	var fromIndex int
	var fromAddress string
	for {
		fromIndexStr := c.readInput("Выберите номер адреса с которого хотите отправить эфир (0-4): ")
		var err error
		fromIndex, err = strconv.Atoi(fromIndexStr)
		if err != nil || fromIndex < 0 || fromIndex >= len(wallet.Addresses) {
			fmt.Println("❌ Неверный номер адреса, попробуйте снова")
			continue
		}
		fromAddress = wallet.Addresses[fromIndex]
		break
	}

	// Цикл для ввода адреса получателя
	var toAddress string
	for {
		toAddress = c.readInput("Введите адрес получателя: ")
		if strings.TrimSpace(toAddress) == "" {
			fmt.Println("❌ Адрес получателя не может быть пустым, попробуйте снова")
			continue
		}
		// Простая проверка формата Ethereum адреса
		if !strings.HasPrefix(toAddress, "0x") || len(toAddress) != 42 {
			fmt.Println("❌ Неверный формат адреса Ethereum (должен начинаться с 0x и содержать 42 символа)")
			continue
		}
		break
	}

	// Цикл для ввода суммы
	var amount string
	for {
		amount = c.readInput("Введите сумму в ETH: ")
		if strings.TrimSpace(amount) == "" {
			fmt.Println("❌ Сумма не может быть пустой, попробуйте снова")
			continue
		}
		// Проверяем, что сумма является числом
		if _, err := strconv.ParseFloat(amount, 64); err != nil {
			fmt.Println("❌ Неверный формат суммы, введите число")
			continue
		}
		break
	}

	fmt.Printf("📋 Подтверждение транзакции:\n")
	fmt.Printf("От: %s\n", fromAddress)
	fmt.Printf("Кому: %s\n", toAddress)
	fmt.Printf("Сумма: %s ETH\n", amount)

	// Цикл для подтверждения
	for {
		confirm := c.readInput("Подтвердить отправку? (y/N): ")
		confirmLower := strings.ToLower(strings.TrimSpace(confirm))

		switch confirmLower {
		case "y", "yes":
			fmt.Println("⏳ Отправка транзакции...")
			result := c.service.SendETH(fromAddress, toAddress, amount)
			fmt.Println(string(result))
			return
		case "n", "no", "":
			fmt.Println("❌ Транзакция отменена")
			return
		default:
			fmt.Println("❌ Введите 'y' для подтверждения или 'n' для отмены")
			continue
		}
	}
}
