package block_getter

import (
	"fmt"
	"sync"

	"github.com/blocto/solana-go-sdk/rpc"
	
	"solana-program-scanner/config"
	"solana-program-scanner/getter/rate_limiter"
	"solana-program-scanner/log"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/parser"
	"solana-program-scanner/sequencer"
)

type BlockGetterManager struct {
	conf           *config.BlockGetterManagerConf
	rpcUrl         string
	startSlot      uint64
	blockParser    *parser.BlockParser
	blockSequencer sequencer.BlockSequencer
	rpcRateLimiter rate_limiter.RpcRateLimiter
	slotBroker     msg_broker.MsgConsumer[uint64]
	getters        []*BlockGetter
}

func NewBlockGetterManager(config *BlockGetterManagerConfig) *BlockGetterManager {
	m := &BlockGetterManager{
		conf:           config.Conf,
		rpcUrl:         config.RPCEndpointHTTP,
		startSlot:      config.StartSlot,
		blockParser:    config.BlockParser,
		blockSequencer: config.BlockSequencer,
		rpcRateLimiter: config.RpcRateLimiter,
		slotBroker:     config.SlotBroker,
	}

	m.createGetters()
	m.initSequencer(m.startSlot)
	return m
}

func (m *BlockGetterManager) createGetters() {
	for i := 0; i < m.conf.GetterNumber; i++ {
		rpcClient := rpc.NewRpcClient(m.rpcUrl)
		id := fmt.Sprintf("BG%d", i)
		getter := NewBlockGetter(
			m.conf.BlockGetter,
			id,
			m.slotBroker,
			&rpcClient,
			m.blockParser,
			m.blockSequencer,
		)
		getter.SubSlotBegin(m.rpcRateLimiter)
		getter.SubSlotErr(m.rpcRateLimiter)
		getter.SubSlotDone(m.rpcRateLimiter)
		m.getters = append(m.getters, getter)
	}
}

func (m *BlockGetterManager) initSequencer(slot uint64) {
	height := m.getters[0].GetBlockHeight(slot)
	if height == 0 {
		log.Logger.Fatal(fmt.Sprintf("StartSlot:%d height=0", config.G.BlockTaskDispatcher.StartSlot))
	}
	m.blockSequencer.Init(height - 1)
}

func (m *BlockGetterManager) Run(wg *sync.WaitGroup) {
	defer func() {
		m.blockSequencer.Close()
		wg.Done()
	}()

	wgWorkers := &sync.WaitGroup{}
	wgWorkers.Add(len(m.getters))
	for _, worker := range m.getters {
		go worker.run(wgWorkers)
	}
	wgWorkers.Wait()
}
