package eth

import "math/big"

type Etherium interface {
	ChainID() (*big.Int, error)
	GetBalance(address string) (*big.Int, error)
	GetTransactionCount(address string) (uint64, error)
	GetGasPrice() (*big.Int, error)
	SendTransaction(tx string) (string, error)
	GetTransactionReceipt(tx string) (string, error)
	GetTransactionByHash(tx string) (string, error)
	SignTransaction(privateKeyHex string, to string, value *big.Int, gasLimit uint64, gasPrice *big.Int, nonce uint64, data []byte, chainID *big.Int) (string, error)
}
