package entities

type Wallet struct {
	privateKeys []string
	Addresses   []string `json:"addresses"` // 5 адресов кошелька
	Seed        []byte   `json:"seed"`      // seed для генерации адресов
}

func (w *Wallet) GetPrivateKeys() []string {
	return w.privateKeys
}

func (w *Wallet) GetPrivateKey(index int) string {
	if index < 0 || index >= len(w.privateKeys) {
		return ""
	}
	return w.privateKeys[index]
}

func (w *Wallet) SetPrivateKeys(keys []string) {
	w.privateKeys = keys
}
