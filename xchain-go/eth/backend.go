package eth

import (
	"math/big"
	"sync"
	"xchain-go/accounts"
	"xchain-go/common"
	"xchain-go/consensus"
	"xchain-go/core"
	"xchain-go/ethdb"
	"xchain-go/event"

	"github.com/ethereum/go-ethereum/miner"
)

// Ethereum implements the Ethereum full node service.
type Ethereum struct {
	// config      *Config
	// chainConfig *params.ChainConfig

	// Channel for shutting down the service
	// shutdownChan chan bool // Channel for shutting down the Ethereum

	// Handlers
	txPool     *core.TxPool
	blockchain *core.BlockChain
	// protocolManager *ProtocolManager
	// lesServer       LesServer

	// DB interfaces
	chainDb ethdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	// bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	// bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *EthAPIBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	etherbase common.Address

	networkId uint64
	// netRPCService *ethapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)
}

func (s *Ethereum) AccountManager() *accounts.Manager { return s.accountManager }
