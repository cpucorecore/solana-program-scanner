package block_getter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"

	"solana-program-scanner/config"
	"solana-program-scanner/getter/pubsub"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/parser"
	"solana-program-scanner/sequencer"
	"solana-program-scanner/types"
)

var (
	transactionVersion = uint8(0)
	rewards            = false
	getBlockConfig     = rpc.GetBlockConfig{
		Encoding:                       rpc.GetBlockConfigEncodingJsonParsed,
		TransactionDetails:             rpc.GetBlockConfigTransactionDetailsFull,
		Rewards:                        &rewards,
		Commitment:                     rpc.Commitment(config.G.Solana.Commitment),
		MaxSupportedTransactionVersion: &transactionVersion,
	}
)

type BlockGetter struct {
	id              string
	slotConsumer    msg_broker.MsgConsumer[uint64]
	rpcClient       *rpc.RpcClient
	ctx             context.Context
	getBlockTimeout time.Duration
	buffer          msg_broker.MsgBroker[*types.BlockWithSlot]
	blockParser     *parser.BlockParser
	blockSequencer  sequencer.BlockSequencer
	slotBeginSub    []pubsub.SlotBeginSub
	slotDoneSub     []pubsub.SlotDoneSub
	slotErrSub      []pubsub.SlotErrSub
}

func NewBlockGetter(
	conf *config.BlockGetterConf,
	id string,
	slotConsumer msg_broker.MsgConsumer[uint64],
	rpcClient *rpc.RpcClient,
	blockParser *parser.BlockParser,
	blockSequencer sequencer.BlockSequencer,
) *BlockGetter {
	return &BlockGetter{
		id:              id,
		slotConsumer:    slotConsumer,
		rpcClient:       rpcClient,
		ctx:             context.Background(),
		getBlockTimeout: time.Millisecond * time.Duration(conf.GetBlockTimeoutByMs),
		buffer:          msg_broker.NewMsgBrokerLocal[*types.BlockWithSlot](conf.BlockBufferSize),
		blockParser:     blockParser,
		blockSequencer:  blockSequencer,
	}
}

func (g *BlockGetter) GetBlockHeight(slot uint64) int64 {
	block, skipped := g.getBlock(slot)
	if skipped {
		return 0
	}

	return *block.BlockHeight
}

func (g *BlockGetter) runBuffer(wg *sync.WaitGroup) {
	defer wg.Done()

	blockBufferCh := g.buffer.ConsumerChan()
	for {
		select {
		case blockWrap, ok := <-blockBufferCh:
			if !ok {
				log.Logger.Info(fmt.Sprintf("%s bufferCh done", g.id))
				return
			}
			log.Logger.Info(fmt.Sprintf("%s slot:%d ParseBlock begin", g.id, blockWrap.Slot))
			parserBlock := g.blockParser.ParseBlock(blockWrap)
			log.Logger.Info(fmt.Sprintf("%s slot:%d ParseBlock finish, Commit begin ", g.id, blockWrap.Slot))
			g.blockSequencer.Commit(parserBlock)
			log.Logger.Info(fmt.Sprintf("%s slot:%d Commit finish", g.id, blockWrap.Slot))
		}
	}
}

func (g *BlockGetter) run(wg *sync.WaitGroup) {
	wgRunBuffer := &sync.WaitGroup{}
	wgRunBuffer.Add(1)
	go g.runBuffer(wgRunBuffer)

	defer func() {
		wgRunBuffer.Wait()
		wg.Done()
	}()

	slotCh := g.slotConsumer.ConsumerChan()
	bufferCh := g.buffer.ProducerChan()
	for {
		select {
		case slot, ok := <-slotCh:
			if !ok {
				log.Logger.Info(fmt.Sprintf("%s slotCh done", g.id))
				g.buffer.Close()
				return
			}

			block, skipped := g.getBlock(slot)
			if skipped {
				log.Logger.Info(fmt.Sprintf("%s slot:%d skipped", g.id, slot))
				continue
			}

			log.Logger.Info(fmt.Sprintf("%s slot:%d commit to buffer begin, buffer.Lenght() %d", g.id, slot, len(bufferCh)))
			bufferCh <- &types.BlockWithSlot{
				Slot:  slot,
				Block: block,
			}
			log.Logger.Info(fmt.Sprintf("%s slot:%d commit to buffer finish, buffer.Lenght() %d", g.id, slot, len(bufferCh)))
		}
	}
}

func shouldRetry(errCode int) bool {
	if types.RpcErrSlotNotAvailable2 == errCode || types.RpcErrSlotNotAvailable == errCode || types.RpcErrSlotCleanup == errCode || types.RpcErrRpsLimit == errCode {
		return true
	}
	return false
}

func slotSkipped(errCode int) bool {
	if types.RpcErrSlotSkippedLedgerJump == errCode || types.RpcErrSlotSkippedLongTerm == errCode {
		return true
	}
	return false
}

var failCnt int

func (g *BlockGetter) GetBlockWithTimeout(slot uint64) (rpc.JsonRpcResponse[*rpc.GetBlock], error) {
	ctx, cancel := context.WithTimeout(g.ctx, g.getBlockTimeout)
	defer cancel()

	return g.rpcClient.GetBlockWithConfig(ctx, slot, getBlockConfig)
}

func monitorBlockDelay(id string, slot uint64, blockTime int64) {
	blockDelay := time.Now().UnixMilli() - (blockTime * 1000)
	monitor.BlockDelay.Observe(float64(blockDelay))
	log.Logger.Info(fmt.Sprintf("%s slot:%d delay %dms", id, slot, blockDelay))
}

func (g *BlockGetter) getBlock(slot uint64) (*rpc.GetBlock, bool) {
	log.Logger.Debug(fmt.Sprintf("%s getBlock:%d start", g.id, slot))

	failCnt = 0
	for {
		g.pubSlotBegin(slot)
		resp, err := g.GetBlockWithTimeout(slot)
		if err != nil {
			failCnt++
			log.Logger.Error(fmt.Sprintf("%s getBlock:%d failCnt:%d with err:%s", g.id, slot, failCnt, err))
			g.pubSlotErr(slot, types.RpcErr)
			continue
		}

		if resp.Error != nil {
			g.pubSlotErr(slot, resp.Error.Code)
			if shouldRetry(resp.Error.Code) {
				failCnt++
				log.Logger.Warn(fmt.Sprintf("%s getBlock:%d failCnt:%d with err:%s, will retry", g.id, slot, failCnt, resp.Error))
				continue
			}

			if slotSkipped(resp.Error.Code) {
				log.Logger.Warn(fmt.Sprintf("%s getBlock:%d skipped with err:%s", g.id, slot, resp.Error.Error()))
				return nil, true
			}

			log.Logger.Fatal(fmt.Sprintf("%s getBlock:%d unknown err:%v", g.id, slot, resp.Error))
			return nil, true
		}

		if resp.Result == nil {
			log.Logger.Warn(fmt.Sprintf("%s getBlock empty:%d err:%s", g.id, slot, resp.Error))
			continue
		}

		g.pubSlotDone(slot)
		monitorBlockDelay(g.id, slot, *resp.Result.BlockTime)

		return resp.Result, false
	}
}

func (g *BlockGetter) SubSlotBegin(sub pubsub.SlotBeginSub) {
	g.slotBeginSub = append(g.slotBeginSub, sub)
}

func (g *BlockGetter) SubSlotDone(sub pubsub.SlotDoneSub) {
	g.slotDoneSub = append(g.slotDoneSub, sub)
}

func (g *BlockGetter) SubSlotErr(sub pubsub.SlotErrSub) {
	g.slotErrSub = append(g.slotErrSub, sub)
}

func (g *BlockGetter) pubSlotBegin(slot uint64) {
	for _, sub := range g.slotBeginSub {
		sub.PubSlotBegin(slot)
	}
}

func (g *BlockGetter) pubSlotDone(slot uint64) {
	for _, sub := range g.slotDoneSub {
		sub.PubSlotDone(slot)
	}
}

func (g *BlockGetter) pubSlotErr(slot uint64, errCode int) {
	for _, sub := range g.slotErrSub {
		sub.PubSlotErr(slot, errCode)
	}
}
