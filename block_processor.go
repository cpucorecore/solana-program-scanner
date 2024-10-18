package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"
	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"go.mongodb.org/mongo-driver/v2/bson"

	"solana-program-scanner/idls/raydium_amm"
)

type BlockProcessor interface {
	name() string
	process(*rpc.GetBlock) error
	done()
}

type BlockProcessorAdmin interface {
	run(ctx context.Context, wg *sync.WaitGroup)
}

type blockProcessorAdmin struct {
	blockCh    chan *rpc.GetBlock
	processors []BlockProcessor
}

func (b *blockProcessorAdmin) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		for _, processor := range b.processors {
			processor.done()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case block := <-b.blockCh:
			if block == nil {
				Logger.Info("blockCh @ done")
				return
			}

			for _, processor := range b.processors {
				err := processor.process(block)
				if err != nil {
					Logger.Error(fmt.Sprintf("%s process block:%d failed", processor.name(), block.BlockHeight))
					return
				}
			}
		}
	}
}

func NewBlockProcessorAdmin(blockCh chan *rpc.GetBlock, txCh chan *OrmTx, ixCh chan bson.M, mg MarketGetter) BlockProcessorAdmin {
	bpf, err := newBlockProcessorFile()
	if err != nil {
		Logger.Fatal(fmt.Sprintf("newBlockProcessorFile err:%v", err))
	}

	rc := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	cm := NewCacheRedisMarket(rc)
	bpp, err := newBlockProcessorParser(txCh, ixCh, mg, cm)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("newBlockProcessorParser err:%v", err))
	}

	processors := []BlockProcessor{
		bpf,
		bpp,
	}

	return &blockProcessorAdmin{
		blockCh:    blockCh,
		processors: processors,
	}
}

type BlockProcessorFile struct {
	fd *os.File
}

var _ BlockProcessor = &BlockProcessorFile{}

const (
	BlocksFilePath = "blocks.json"
)

func newBlockProcessorFile() (bpf *BlockProcessorFile, err error) {
	f, err := os.Create(BlocksFilePath)
	if err != nil {
		Logger.Error(fmt.Sprintf("create file err:%s", err.Error()))
		return nil, err
	}
	return &BlockProcessorFile{
		fd: f,
	}, nil
}

func (bpf *BlockProcessorFile) name() string {
	return "file"
}

func (bpf *BlockProcessorFile) process(block *rpc.GetBlock) error {
	blockData, err := json.Marshal(block)
	if err != nil {
		Logger.Error(fmt.Sprintf(""))
		return err
	}

	_, err = bpf.fd.Write(blockData)
	return err
}

func (bpf *BlockProcessorFile) done() {
	bpf.fd.Close()
}

type BlockProcessorParser struct {
	txCh         chan *OrmTx
	ixCh         chan bson.M
	marketGetter MarketGetter
	cacheMarket  Cache[string, *Market]
}

func newBlockProcessorParser(txCh chan *OrmTx, ixCh chan bson.M, mg MarketGetter, cm Cache[string, *Market]) (bpp *BlockProcessorParser, err error) {
	return &BlockProcessorParser{
		ixCh:         ixCh,
		txCh:         txCh,
		marketGetter: mg,
		cacheMarket:  cm,
	}, nil
}

func (bpp *BlockProcessorParser) name() string {
	return "parser"
}

const (
	RadiumAmmAddressMainnet = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

func (bpp *BlockProcessorParser) process(block *rpc.GetBlock) error {
	for txIdx, tx := range block.Transactions {
		if tx.Meta.Err != nil {
			continue
		}

		accountKeys := make(map[string]rpc.AccountKey)
		for _, accountKey := range tx.Transaction.Message.AccountKeys {
			accountKeys[accountKey.Pubkey] = accountKey
		}

		var ixF rpc.InstructionFull
		for _, instruction := range tx.Transaction.Message.Instructions {
			ixJson, err := json.Marshal(instruction)
			if err != nil {
				Logger.Fatal(fmt.Sprintf("json marshal err:%s", err.Error()))
			}

			err = json.Unmarshal(ixJson, &ixF)
			if err != nil {
				Logger.Fatal(fmt.Sprintf("json unmarshal err:%s", err.Error()))
			}

			if ixF.ProgramId == RadiumAmmAddressMainnet {
				var accounts []*ag_solanago.AccountMeta
				for _, account := range ixF.Accounts {
					accounts = append(accounts, &ag_solanago.AccountMeta{
						PublicKey:  ag_solanago.MustPublicKeyFromBase58(account),
						IsWritable: accountKeys[account].Writable,
						IsSigner:   accountKeys[account].Signer,
					})
				}

				ixData, err := base58.Decode(ixF.Data)
				if err != nil {
					Logger.Error(fmt.Sprintf("decode err:%s", err.Error()))
				}

				ix, err := raydium_amm.DecodeInstruction(accounts, ixData)
				if err != nil {
					Logger.Error(fmt.Sprintf("decode instruction err:%s", err.Error()))
				}

				switch ix.TypeID.Uint8() {
				case raydium_amm.Instruction_SwapBaseIn:
					ixNude, ok := ix.Impl.(*raydium_amm.SwapBaseIn)
					if ok {
						marketAddress := ixNude.AccountMetaSlice[1].PublicKey.String()
						market, err := bpp.cacheMarket.Get(marketAddress)
						if err != nil {
							Logger.Fatal(fmt.Sprintf("market cache get err:%v", err)) // TODO check
						}

						if market == nil {
							market, err = bpp.marketGetter.getMarket(marketAddress) // TODO save to db table:market/pair
							if err != nil {
								Logger.Fatal(fmt.Sprintf("get market err:%v", err))
							}

							err = bpp.cacheMarket.Set(market.Address, market)
							if err != nil {
								Logger.Warn(fmt.Sprintf("market cache set err:%v", err))
							}
						}

						ormTx := OrmTx{
							TxHash:        tx.Transaction.Signatures[0],
							Event:         int8(len(tx.Meta.LogMessages)), // TODO check
							Token0Amount:  strconv.FormatUint(*ixNude.AmountIn, 10),
							Token1Amount:  strconv.FormatUint(*ixNude.MinimumAmountOut, 10), // TODO check 0
							Maker:         ixNude.AccountMetaSlice[1].PublicKey.String(),    // TODO
							Token0Address: market.BaseMint,
							Token1Address: market.QuoteMint,
							Block:         *block.BlockHeight,
							BlockAt:       time.Unix(*block.BlockTime, 0),
							Index:         txIdx,
						}

						txJson, _ := json.Marshal(tx)
						Logger.Debug(fmt.Sprintf("tx:%s", string(txJson)))

						bpp.txCh <- &ormTx

						ixJson, err = json.Marshal(ixNude)
						if err != nil {
							Logger.Error(fmt.Sprintf("json marshal err:%s", err.Error()))
						}

						Logger.Debug(fmt.Sprintf("ix:%s", string(ixJson)))
						bpp.ixCh <- bson.M{"data": string(ixJson)}
					} else {
						Logger.Error(fmt.Sprintf("instruction:%s assertion failed", raydium_amm.InstructionIDToName(ix.TypeID.Uint8())))
					}
					break
				case raydium_amm.Instruction_SwapBaseOut:
					ixNude, ok := ix.Impl.(*raydium_amm.SwapBaseOut)
					if ok {
						ixJson, err = json.Marshal(ixNude)
						if err != nil {
							Logger.Error(fmt.Sprintf("json marshal err:%s", err.Error()))
						}

						bpp.ixCh <- bson.M{"data": string(ixJson)}
					} else {
						Logger.Error(fmt.Sprintf("instruction:%s assertion failed", raydium_amm.InstructionIDToName(ix.TypeID.Uint8())))
					}
					break
				default:
				}
			}
		}
	}

	return nil
}

func (bpp *BlockProcessorParser) done() {
	close(bpp.ixCh)
	close(bpp.txCh)
}

var _ BlockProcessor = &BlockProcessorParser{}
