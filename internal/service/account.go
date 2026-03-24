package service

import (
	"fmt"
	"strings"

	"github.com/Murolando/ftem_wallet/pkg/entities"
	"github.com/Murolando/ftem_wallet/pkg/lib/algorithm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
)

// AuthWallet авторизует кошелек
func (s *Service) AuthWallet(mnemonicWords [12]string, password string) entities.ResultString {
	// Проверяем валидность мнемонической фразы
	if !algorithm.BIP39IsValidMnemomic(mnemonicWords) {
		return entities.ErrInvalidMnemonic
	}

	// Получаем seed из мнемоники
	seed := algorithm.BIP39SeedFromMnemomic(mnemonicWords, password)

	// Заполняем кошелек
	err := s.getWallet(seed)
	if err != nil {
		return entities.ErrPrivateKeyDerivation
	}

	result := "Wallet successfully authorized!\n\nAddresses:\n"
	for i, addr := range s.currentWallet.Addresses {
		result += fmt.Sprintf("Wallet %d: %s\n", i, addr)
	}

	return entities.ResultString(result)
}

// GenerateWallet генерирует новый кошелек
func (s *Service) GenerateWallet(password string) entities.ResultString {
	// Генерируем мнемоническую фразу (12 слов)
	mnemonicWords := algorithm.BIP39Mnemonic()
	mnemonicPhrase := strings.Join(mnemonicWords[:], " ")

	// Получаем seed из мнемоники
	seed := algorithm.BIP39SeedFromMnemomic(mnemonicWords, password)

	// Заполняем кошелек
	err := s.getWallet(seed)
	if err != nil {
		return entities.ErrPrivateKeyDerivation
	}

	result := fmt.Sprintf("Mnemonic: %s\n\nAddresses:\n", mnemonicPhrase)
	for i, addr := range s.currentWallet.Addresses {
		result += fmt.Sprintf("Wallet %d: %s\n", i, addr)
	}

	return entities.ResultString(result)
}

func (s *Service) getWallet(seed []byte) error {
	// Получаем 5 адресов и приватных ключей
	addresses := make([]string, 0, 5)
	privateKeys := make([]string, 0, 5)

	// Деривация по пути BIP44: m/44'/60'/0'/0/i
	master := algorithm.BIP32Master(seed)
	purpose := algorithm.BIP32Child(master, bip32.FirstHardenedChild+44)
	coinType := algorithm.BIP32Child(purpose, bip32.FirstHardenedChild+60)
	account := algorithm.BIP32Child(coinType, bip32.FirstHardenedChild+0)
	change := algorithm.BIP32Child(account, 0)

	for i := uint32(0); i < 5; i++ {
		child := algorithm.BIP32Child(change, i)

		// Получаем приватный ключ
		priv, err := crypto.ToECDSA(child.Key)
		if err != nil {
			return err
		}

		// Сохраняем приватный ключ в hex формате
		privateKey := fmt.Sprintf("%x", crypto.FromECDSA(priv))
		privateKeys = append(privateKeys, privateKey)

		// Получаем адрес кошелька
		addr := crypto.PubkeyToAddress(priv.PublicKey).Hex()
		addresses = append(addresses, addr)
	}

	wallet := &entities.Wallet{
		Addresses: addresses,
		Seed:      seed,
	}
	wallet.SetPrivateKeys(privateKeys)
	s.SetCurrentWallet(wallet)

	return nil
}
