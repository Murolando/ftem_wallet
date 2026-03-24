# 🔐 FTEM Wallet

*[English version below](#english-version)*

## 📖 Описание

FTEM Wallet - это консольное приложение для управления Ethereum кошельком, написанное на Go. Приложение поддерживает создание кошельков с использованием мнемонических фраз BIP-39, проверку баланса и отправку ETH транзакций.

## ✨ Возможности

- 🆕 **Создание нового кошелька** с генерацией мнемонической фразы (12 слов)
- 🔑 **Импорт кошелька** из существующей мнемонической фразы
- 👁️ **Просмотр адресов** кошелька (поддержка до 5 адресов)
- 💰 **Проверка баланса** всех адресов кошелька
- 📤 **Отправка ETH** между адресами
- 🔒 **Безопасное хранение** с поддержкой паролей
- 🌐 **Поддержка тестовых сетей** (Sepolia)

## 🏗️ Архитектура

Проект следует принципам чистой архитектуры:

```
ftem_wallet/
├── cmd/                    # Точки входа приложения
│   ├── dev/               # Версия для разработки
│   └── prod/              # Продакшн версия
├── internal/              # Внутренняя логика приложения
│   ├── cli/               # CLI интерфейс
│   └── service/           # Бизнес-логика
├── pkg/                   # Публичные пакеты
│   ├── config/            # Конфигурация
│   ├── entities/          # Сущности
│   └── lib/               # Библиотеки
│       ├── algorithm/     # Криптографические алгоритмы
│       └── clients/       # Клиенты для внешних API
├── configs/               # Файлы конфигурации
└── releases/              # Скомпилированные бинарные файлы
```

## 🛠️ Технологии

- **Go 1.25.5** - основной язык программирования
- **go-ethereum** - взаимодействие с Ethereum блокчейном
- **BIP-32/BIP-39** - стандарты для генерации кошельков
- **Viper** - управление конфигурацией
- **Infura** - провайдер для подключения к Ethereum сети

## 📦 Установка

### Предварительные требования

- Go 1.25.5 или выше
- Git

### Клонирование репозитория

```bash
git clone https://github.com/Murolando/ftem_wallet.git
cd ftem_wallet
```

## 🚀 Запуск

### Режим разработки (Sepolia Testnet)

```bash
# Запуск с конфигурацией для разработки
go run cmd/dev/main.go --app-cfg configs/dev/config.yaml

# Или с переменной окружения
export APP_CFG_PATH=configs/dev/config.yaml
go run cmd/dev/main.go

# С мнемонической фразой
go run cmd/dev/main.go --app-cfg configs/dev/config.yaml --mnemonic "your twelve word mnemonic phrase here goes like this example phrase" --password "your_password"
```

### Продакшн режим (Ethereum Mainnet)

#### Способ 1: Прямой запуск
```bash
# Запуск prod версии с конфигурацией
go run cmd/prod/main.go --app-cfg configs/prod/config.yaml

# Или с переменной окружения
export APP_CFG_PATH=configs/prod/config.yaml
go run cmd/prod/main.go

# С мнемонической фразой
go run cmd/prod/main.go --app-cfg configs/prod/config.yaml --mnemonic "your mnemonic phrase" --password "your_password"
```

#### Способ 2: Сборка и запуск бинарного файла
```bash
# Сборка prod версии
go build -o releases/ftem_wallet_prod cmd/prod/main.go

# Запуск с конфигурацией
./releases/ftem_wallet_prod --app-cfg configs/prod/config.yaml

# Или с переменной окружения
export APP_CFG_PATH=configs/prod/config.yaml
./releases/ftem_wallet_prod

# С мнемонической фразой
./releases/ftem_wallet_prod --app-cfg configs/prod/config.yaml --mnemonic "your mnemonic phrase" --password "your_password"
```

### ⚠️ Важные замечания для продакшна

1. **Обязательно замените** `YOUR_PROJECT_ID` в [`configs/prod/config.yaml`](configs/prod/config.yaml) на ваш реальный Infura Project ID
2. **Используйте реальные ETH** - prod версия работает с Ethereum Mainnet
3. **Сохраните мнемоническую фразу** в безопасном месте
4. **Проверьте баланс** перед отправкой транзакций

## ⚙️ Конфигурация

Приложение использует YAML файлы конфигурации:

### configs/dev/config.yaml
```yaml
env: dev
service_config:
  host: https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID
cli_config:
  mnemonic: ""
  password: ""
```

### configs/prod/config.yaml
```yaml
env: prod
service_config:
  host: https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID
cli_config:
  mnemonic: ""
  password: ""
```

## 📋 Использование

### Создание нового кошелька

1. Запустите приложение
2. Выберите опцию "1. 🆕 Создать новый кошелек"
3. Введите пароль (опционально)
4. **ВАЖНО**: Сохраните сгенерированную мнемоническую фразу в безопасном месте

### Импорт существующего кошелька

Запустите приложение с флагами:
```bash
go run cmd/dev/main.go --mnemonic "your mnemonic phrase" --password "your_password"
```

### Основные операции

После авторизации доступны следующие операции:

1. **👁️ Показать адреса кошелька** - отображает все адреса (до 5)
2. **💰 Показать баланс** - проверяет баланс всех адресов
3. **📤 Отправить ETH** - отправка ETH между адресами
4. **❌ Выход** - завершение работы

## 🔒 Безопасность

- Мнемонические фразы генерируются с использованием криптографически стойкого генератора случайных чисел
- Поддержка паролей для дополнительной защиты
- Приватные ключи не сохраняются в открытом виде
- Все операции требуют подтверждения пользователя

## 🌐 Поддерживаемые сети

- **Sepolia Testnet** (режим разработки)
- **Ethereum Mainnet** (продакшн режим)

---

# English Version

## 📖 Description

FTEM Wallet is a command-line application for managing Ethereum wallets, written in Go. The application supports wallet creation using BIP-39 mnemonic phrases, balance checking, and sending ETH transactions.

## ✨ Features

- 🆕 **Create new wallet** with mnemonic phrase generation (12 words)
- 🔑 **Import wallet** from existing mnemonic phrase
- 👁️ **View wallet addresses** (supports up to 5 addresses)
- 💰 **Check balance** of all wallet addresses
- 📤 **Send ETH** between addresses
- 🔒 **Secure storage** with password support
- 🌐 **Testnet support** (Sepolia)

## 🏗️ Architecture

The project follows clean architecture principles:

```
ftem_wallet/
├── cmd/                    # Application entry points
│   ├── dev/               # Development version
│   └── prod/              # Production version
├── internal/              # Internal application logic
│   ├── cli/               # CLI interface
│   └── service/           # Business logic
├── pkg/                   # Public packages
│   ├── config/            # Configuration
│   ├── entities/          # Entities
│   └── lib/               # Libraries
│       ├── algorithm/     # Cryptographic algorithms
│       └── clients/       # External API clients
├── configs/               # Configuration files
└── releases/              # Compiled binaries
```

## 📦 Installation

### Prerequisites

- Go 1.25.5 or higher
- Git

### Clone the repository

```bash
git clone https://github.com/Murolando/ftem_wallet.git
cd ftem_wallet
```

## 🚀 Running

### Development mode (Sepolia Testnet)

```bash
# Run with development configuration
go run cmd/dev/main.go --app-cfg configs/dev/config.yaml

# Or with environment variable
export APP_CFG_PATH=configs/dev/config.yaml
go run cmd/dev/main.go

# With mnemonic phrase
go run cmd/dev/main.go --app-cfg configs/dev/config.yaml --mnemonic "your twelve word mnemonic phrase here goes like this example phrase" --password "your_password"
```

### Production mode (Ethereum Mainnet)

#### Method 1: Direct run
```bash
# Run prod version with configuration
go run cmd/prod/main.go --app-cfg configs/prod/config.yaml

# Or with environment variable
export APP_CFG_PATH=configs/prod/config.yaml
go run cmd/prod/main.go

# With mnemonic phrase
go run cmd/prod/main.go --app-cfg configs/prod/config.yaml --mnemonic "your mnemonic phrase" --password "your_password"
```

#### Method 2: Build and run binary
```bash
# Build prod version
go build -o releases/ftem_wallet_prod cmd/prod/main.go

# Run with configuration
./releases/ftem_wallet_prod --app-cfg configs/prod/config.yaml

# Or with environment variable
export APP_CFG_PATH=configs/prod/config.yaml
./releases/ftem_wallet_prod

# With mnemonic phrase
./releases/ftem_wallet_prod --app-cfg configs/prod/config.yaml --mnemonic "your mnemonic phrase" --password "your_password"
```

### ⚠️ Important notes for production

1. **Must replace** `YOUR_PROJECT_ID` in [`configs/prod/config.yaml`](configs/prod/config.yaml) with your real Infura Project ID
2. **Uses real ETH** - prod version works with Ethereum Mainnet
3. **Save mnemonic phrase** in a secure place
4. **Check balance** before sending transactions

## ⚙️ Configuration

The application uses YAML configuration files:

### configs/dev/config.yaml
```yaml
env: dev
service_config:
  host: https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID
cli_config:
  mnemonic: ""
  password: ""
```

### configs/prod/config.yaml
```yaml
env: prod
service_config:
  host: https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID
cli_config:
  mnemonic: ""
  password: ""
```

## 📋 Usage

### Creating a new wallet

1. Run the application
2. Select option "1. 🆕 Create new wallet"
3. Enter password (optional)
4. **IMPORTANT**: Save the generated mnemonic phrase in a secure place

### Import existing wallet

Run the application with flags:
```bash
go run cmd/dev/main.go --mnemonic "your mnemonic phrase" --password "your_password"
```

### Main operations

After authorization, the following operations are available:

1. **👁️ Show wallet addresses** - displays all addresses (up to 5)
2. **💰 Show balance** - checks balance of all addresses
3. **📤 Send ETH** - send ETH between addresses
4. **❌ Exit** - terminate the application

## 🔒 Security

- Mnemonic phrases are generated using cryptographically secure random number generator
- Password support for additional protection
- Private keys are not stored in plain text
- All operations require user confirmation

## 🌐 Supported Networks

- **Sepolia Testnet** (development mode)
- **Ethereum Mainnet** (production mode)
