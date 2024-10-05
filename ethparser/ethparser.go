package ethparser

import (
	"sync"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
	// PollBlockchain check the eth netowkr for transactions
	PollBlockchain() error
}

type EthParser struct {
	client *ethclient

	mu sync.RWMutex // To handle concurrent access

	currentBlockHex string
	subscribers     map[string]bool
	transactions    map[string][]Transaction
}

// NewEthParser initializes the parser and start polling the blockchain
func NewEthParser() Parser {
	client := NewEthClient("")
	parser := &EthParser{
		client:          client,
		currentBlockHex: "0",
		subscribers:     make(map[string]bool),
		transactions:    make(map[string][]Transaction),
	}

	return parser
}

// GetCurrentBlock returns the last parsed block number
func (p *EthParser) GetCurrentBlock() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	block, err := hexStringToInt(p.currentBlockHex)
	if err != nil {
		return 0
	}

	return block
}

// Subscribe adds an address to the observer list
func (p *EthParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.subscribers[address]; exists {
		return false
	}
	p.subscribers[address] = true
	return true
}

// GetTransactions returns a list of inbound or outbound transactions for a subscribed address
func (p *EthParser) GetTransactions(address string) []Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.transactions[address]
}

// fetchTransactionsForBlock fetches all the inbound or outbound transactions for a block from ethereum network
func (p *EthParser) fetchTransactionsForBlock(blockHex string) error {
	allTxns, err := p.client.GetBlockTransactions(blockHex)
	if err != nil {
		return err
	}

	for _, tx := range allTxns {
		p.mu.Lock()

		// check inbound transaction
		if _, ok := p.subscribers[tx.To]; ok {
			p.transactions[tx.To] = append(p.transactions[tx.To], tx)
		}

		// check outbound transaction
		if _, ok := p.subscribers[tx.From]; ok {
			p.transactions[tx.From] = append(p.transactions[tx.From], tx)
		}

		p.mu.Unlock()
	}

	return nil
}

// PollBlockchain keeps on polling blockchian to fetch latest transactions
func (p *EthParser) PollBlockchain() error {
	latestBlockHex, err := p.client.GetLatestBlock()
	if err != nil {
		return err
	}

	latestBlockInt, err := hexStringToInt(latestBlockHex)
	if err != nil {
		return err
	}

	if latestBlockInt > p.GetCurrentBlock() {
		p.mu.Lock()
		p.currentBlockHex = latestBlockHex
		p.mu.Unlock()

		err = p.fetchTransactionsForBlock(latestBlockHex)
		return err
	}

	return nil
}
