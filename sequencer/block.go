package sequencer

import (
	"solana-program-scanner/config"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/parser"
	"sync"
)

type BlockSequencer interface {
	Init(height int64)
	Commit(block *parser.Block)
	Close()
}

type blockSequencer struct {
	active        bool
	mu            sync.Mutex
	cond          *sync.Cond
	height        int64
	blockProducer msg_broker.MsgProducer[*parser.Block]
}

func NewBlockSequencer(conf *config.BlockSequencerConf, blockProducer msg_broker.MsgProducer[*parser.Block]) BlockSequencer {
	s := &blockSequencer{
		active:        conf.Active,
		blockProducer: blockProducer,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *blockSequencer) Init(height int64) {
	s.height = height
}

func (s *blockSequencer) Commit(block *parser.Block) {
	if !s.active {
		s.blockProducer.Produce(block)
		return
	}

	s.mu.Lock()
	for s.height+1 != block.Block.Height {
		s.cond.Wait()
	}

	s.blockProducer.Produce(block)
	s.height = block.Block.Height
	s.cond.Broadcast()
	s.mu.Unlock()
}

func (s *blockSequencer) Close() {
	s.blockProducer.Close()
}
