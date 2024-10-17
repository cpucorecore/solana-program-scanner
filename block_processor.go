package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

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

func NewBlockProcessorAdmin(blockCh chan *rpc.GetBlock, txRawChan chan string, ixRawChan chan string, ixIndexCh chan bson.M, ixCh chan bson.M) BlockProcessorAdmin {
	bpf, err := newBlockProcessorFile()
	if err != nil {
		Logger.Fatal(fmt.Sprintf("newBlockProcessorFile err:%v", err))
	}

	bpp, err := newBlockProcessorParser(txRawChan, ixRawChan, ixIndexCh, ixCh)
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
	txRawChan chan string
	ixRawChan chan string
	ixIndexCh chan bson.M
	ixCh      chan bson.M
}

func newBlockProcessorParser(txRawChan chan string, ixRawChan chan string, ixIndexCh chan bson.M, ixCh chan bson.M) (bpp *BlockProcessorParser, err error) {
	return &BlockProcessorParser{
		txRawChan: txRawChan,
		ixRawChan: ixRawChan,
		ixIndexCh: ixIndexCh,
		ixCh:      ixCh,
	}, nil
}

func (bpp *BlockProcessorParser) name() string {
	return "parser"
}

const (
	OpenbookV2AddressMainnet = "opnb2LAfJYbRMAHHvqjCwQxanZn7ReEHp1k81EohpZb"
	RadiumAmmAddressMainnet  = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

func (bpp *BlockProcessorParser) process(block *rpc.GetBlock) error {
	for _, tx := range block.Transactions {
		if tx.Meta.Err != nil {
			//Logger.Info(fmt.Sprintf("skip failed tx Signatures:%v", tx.Transaction.Signatures))
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
				Logger.Error(fmt.Sprintf("json marshal err:%s", err.Error()))
				// TODO exit
			}

			err = json.Unmarshal(ixJson, &ixF)
			if err != nil {
				Logger.Error(fmt.Sprintf("json unmarshal err:%s", err.Error()))
				// TODO exit
			}

			if ixF.ProgramId == RadiumAmmAddressMainnet {
				txJson, err := json.Marshal(tx)
				if err != nil {
					Logger.Error(fmt.Sprintf("json unmarshal err:%s", err.Error()))
					// TODO exit
				}
				//Logger.Info(string(txJson))
				//Logger.Info(string(ixJson))
				bpp.txRawChan <- string(txJson)
				bpp.ixRawChan <- string(ixJson)

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
				//Logger.Info(fmt.Sprintf("%s", hex.EncodeToString(ixData)))

				ix, err := raydium_amm.DecodeInstruction(accounts, ixData)
				if err != nil {
					Logger.Error(fmt.Sprintf("decode instruction err:%s", err.Error()))
				}

				bpp.ixIndexCh <- bson.M{"ins": raydium_amm.InstructionIDToName(ix.TypeID.Uint8()), "signature": tx.Transaction.Signatures}

				switch ix.TypeID.Uint8() {
				case raydium_amm.Instruction_SwapBaseIn:
					ixNude, ok := ix.Impl.(*raydium_amm.SwapBaseIn)
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
	close(bpp.txRawChan)
	close(bpp.ixRawChan)
	close(bpp.ixIndexCh)
	close(bpp.ixCh)
}

var _ BlockProcessor = &BlockProcessorParser{}
