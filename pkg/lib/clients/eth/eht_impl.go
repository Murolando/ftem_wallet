package eth

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type EtheriumClient struct {
	url string
}

func NewEtheriumClient(url string) *EtheriumClient {
	return &EtheriumClient{
		url: url,
	}
}

func (c *EtheriumClient) ChainID() (*big.Int, error) {
	response, err := c.do(map[string][]string{
		"method": {"eth_chainId"},
	})
	if err != nil {
		return nil, err
	}
	
	return c.parseHexToBigInt(response)
}
// https://cloud.google.com/application/web3/faucet/ethereum/sepolia пополняют бесплатно, спасибо ребята
func (c *EtheriumClient) GetBalance(address string) (*big.Int, error) {
	response, err := c.do(map[string][]string{
		"method": {"eth_getBalance"},
		"params": {address, "latest"},
	})
	if err != nil {
		return nil, err
	}

	return c.parseHexToBigInt(response)
}

func (c *EtheriumClient) GetTransactionCount(address string) (uint64, error) {
	response, err := c.do(map[string][]string{
		"method": {"eth_getTransactionCount"},
		"params": {address, "pending"},
	})
	if err != nil {
		return 0, err
	}

	return c.parseHexToUint64(response)
}

func (c *EtheriumClient) GetGasPrice() (*big.Int, error) {
	response, err := c.do(map[string][]string{
		"method": {"eth_gasPrice"},
	})
	if err != nil {
		return nil, err
	}

	return c.parseHexToBigInt(response)
}

func (c *EtheriumClient) SendTransaction(tx string) (string, error) {
	response, err := c.do(map[string][]string{
		"method": {"eth_sendRawTransaction"},
		"params": {tx},
	})
	if err != nil {
		return "", err
	}

	// Парсим ответ для проверки на ошибки RPC
	var rpcResponse struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal([]byte(response), &rpcResponse); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if rpcResponse.Error != nil {
		return "", fmt.Errorf("RPC error: %s", rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}

func (c *EtheriumClient) GetTransactionReceipt(tx string) (string, error) {
	return c.do(map[string][]string{
		"method": {"eth_getTransactionByHash"},
		"params": {tx},
	})
}

func (c *EtheriumClient) GetTransactionByHash(tx string) (string, error) {
	return c.do(map[string][]string{
		"method": {"eth_getTransactionReceipt"},
		"params": {tx},
	})
}

func (c *EtheriumClient) do(methodParams map[string][]string) (string, error) {
	// Создаем JSON-RPC запрос
	jsonRPCReq := map[string]any{
		"jsonrpc": "2.0",
		"method":  methodParams["method"][0], // берем первый элемент из слайса
		"params":  methodParams["params"],
		"id":      1,
	}

	reqBytes, err := json.Marshal(jsonRPCReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON-RPC request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(respBody), nil
}

// parseHexToBigInt парсит JSON ответ и конвертирует hex значение в big.Int
func (c *EtheriumClient) parseHexToBigInt(jsonResponse string) (*big.Int, error) {
	var response struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal([]byte(jsonResponse), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", response.Error.Message)
	}

	// Конвертируем hex в big.Int
	result := new(big.Int)
	if strings.HasPrefix(response.Result, "0x") {
		result.SetString(response.Result[2:], 16)
	} else {
		result.SetString(response.Result, 16)
	}

	return result, nil
}

// SignTransaction создает и подписывает транзакцию
func (c *EtheriumClient) SignTransaction(privateKeyHex string, to string, value *big.Int, gasLimit uint64, gasPrice *big.Int, nonce uint64, data []byte, chainID *big.Int) (string, error) {
	// 1. Парсим приватный ключ
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// 2. Создаем транзакцию
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(to),
		value,
		gasLimit,
		gasPrice,
		data,
	)

	// 3. Подписываем транзакцию
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 4. Кодируем в RLP для отправки
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction: %w", err)
	}

	return "0x" + hex.EncodeToString(rawTx), nil
}

// parseHexToUint64 парсит JSON ответ и конвертирует hex значение в uint64
func (c *EtheriumClient) parseHexToUint64(jsonResponse string) (uint64, error) {
	var response struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal([]byte(jsonResponse), &response); err != nil {
		return 0, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if response.Error != nil {
		return 0, fmt.Errorf("RPC error: %s", response.Error.Message)
	}

	// Конвертируем hex в uint64
	var result uint64
	var err error
	if strings.HasPrefix(response.Result, "0x") {
		result, err = strconv.ParseUint(response.Result[2:], 16, 64)
	} else {
		result, err = strconv.ParseUint(response.Result, 16, 64)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to convert hex to uint64: %w", err)
	}

	return result, nil
}
