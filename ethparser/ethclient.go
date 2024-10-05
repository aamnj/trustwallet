package ethparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultEthRPCUrl = "https://ethereum-rpc.publicnode.com"

type ethclient struct {
	ethRPCUrl string
}

// NewEthClient returns instance of ethclient to interact with ethereum network
func NewEthClient(url string) *ethclient {
	if url == "" {
		url = defaultEthRPCUrl
	}

	return &ethclient{ethRPCUrl: url}
}

// JSONRPCRequest represents a basic Ethereum JSON-RPC request
type JSONRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSONRPCResponse represents a basic Ethereum JSON-RPC response
type JSONRPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      int         `json:"id"`
}

// GetLatestBlock fetches the latest block number from the Ethereum network
func (c *ethclient) GetLatestBlock() (string, error) {
	response, err := c.makePostRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		return "", err
	}

	blockNumberHex := response.Result.(string)
	return strings.TrimPrefix(blockNumberHex, "0x"), nil
}

// GetBlockTransactions fetches transaction of the current block
func (c *ethclient) GetBlockTransactions(blockHex string) ([]Transaction, error) {
	txns := []Transaction{}
	block := fmt.Sprintf("0x%v", blockHex)
	response, err := c.makePostRequest("eth_getBlockByNumber", []interface{}{block, true})
	if err != nil {
		return txns, err
	}

	res, ok1 := response.Result.(map[string]interface{})
	if !ok1 {
		return txns, err
	}

	transactions, ok2 := res["transactions"].([]interface{})
	if !ok2 {
		return txns, err
	}

	// Debug stmt to track count of transaction in each block
	log.Printf("total transactions in block %v: %v", block, len(transactions))

	for _, txn := range transactions {
		m, ok := txn.(map[string]interface{})
		if ok {
			hash, ok1 := m["hash"].(string)
			from, ok2 := m["from"].(string)
			to, ok3 := m["to"].(string)
			value, ok4 := m["value"].(string)

			if ok1 && ok2 && ok3 && ok4 {
				newTxn := Transaction{
					Hash:  hash,
					From:  from,
					To:    to,
					Value: value,
				}
				txns = append(txns, newTxn)
			}
		}
	}
	return txns, nil
}

func (c *ethclient) makePostRequest(method string, params []interface{}) (JSONRPCResponse, error) {
	var rpcResponse JSONRPCResponse

	requestBody := JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	body, _ := json.Marshal(requestBody)
	resp, err := http.Post(c.ethRPCUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return rpcResponse, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&rpcResponse)
	if err != nil {
		return rpcResponse, err
	}

	return rpcResponse, nil
}
