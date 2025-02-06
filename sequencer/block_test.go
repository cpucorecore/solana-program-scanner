package sequencer

import (
	"sync"
	"testing"
	"time"

	"solana-program-scanner/config"
	"solana-program-scanner/parser"

	"github.com/stretchr/testify/assert"
)

// 模拟消息生产者
type mockBlockProducer struct {
	blocks []*parser.Block
	ch     chan *parser.Block
}

func newMockBlockProducer() *mockBlockProducer {
	return &mockBlockProducer{
		blocks: make([]*parser.Block, 0, 100),
		ch:     make(chan *parser.Block, 100),
	}
}

func (p *mockBlockProducer) ProducerChan() chan<- *parser.Block {
	return p.ch
}

func (p *mockBlockProducer) Produce(block *parser.Block) {
	p.blocks = append(p.blocks, block)
	p.ch <- block
}

func (p *mockBlockProducer) Close() {}

func TestBlockSequencer_Commit(t *testing.T) {
	// 初始化配置
	conf := &config.BlockSequencerConf{
		Active: true,
	}

	// 初始化mock producer
	producer := newMockBlockProducer()

	// 创建sequencer
	sequencer := NewBlockSequencer(conf, producer)
	sequencer.Init(100) // 从高度100开始

	// 创建测试区块
	blocks := make([]*parser.Block, 100)
	for i := 0; i < 100; i++ {
		blocks[i] = &parser.Block{
			Block: &parser.BlockHeader{
				Height: int64(101 + i), // 101, 102, 103, 104, 105
			},
		}
	}

	// 并发提交区块
	var wg sync.WaitGroup
	wg.Add(100)

	// 乱序提交区块

	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			sequencer.Commit(blocks[i])
		}(i)
	}

	// 等待所有区块提交完成
	wg.Wait()

	// 验证提交顺序
	assert.Equal(t, 100, len(producer.blocks))
	for i := 0; i < 100; i++ {
		assert.Equal(t, int64(101+i), producer.blocks[i].Block.Height)
	}
}

func TestBlockSequencer_Inactive(t *testing.T) {
	// 测试inactive状态
	conf := &config.BlockSequencerConf{
		Active: false,
	}
	producer := newMockBlockProducer()
	sequencer := NewBlockSequencer(conf, producer)

	block := &parser.Block{
		Block: &parser.BlockHeader{
			Height: 101,
		},
	}

	sequencer.Commit(block)
	assert.Equal(t, 1, len(producer.blocks))
}

func TestBlockSequencer_NonSequentialHeight(t *testing.T) {
	// 测试不连续的区块高度
	conf := &config.BlockSequencerConf{
		Active: true,
	}
	producer := newMockBlockProducer()
	sequencer := NewBlockSequencer(conf, producer)
	sequencer.Init(100)

	var wg sync.WaitGroup
	wg.Add(2)

	// 提交高度为102的区块
	go func() {
		defer wg.Done()
		sequencer.Commit(&parser.Block{
			Block: &parser.BlockHeader{Height: 102},
		})
	}()

	// 等待一小段时间后提交高度为101的区块
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		sequencer.Commit(&parser.Block{
			Block: &parser.BlockHeader{Height: 101},
		})
	}()

	wg.Wait()

	assert.Equal(t, 2, len(producer.blocks))
	assert.Equal(t, int64(101), producer.blocks[0].Block.Height)
	assert.Equal(t, int64(102), producer.blocks[1].Block.Height)
}
