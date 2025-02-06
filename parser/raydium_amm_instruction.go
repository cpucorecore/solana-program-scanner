package parser

import (
	"github.com/blocto/solana-go-sdk/rpc"
	"solana-program-scanner/parser/filters"
)

const InstructionIndexBase = 10000

type RaydiumInstructionWrap struct {
	Raydium   *Instruction
	Transfers []*InstructionParsed
}

type RaydiumInstructions struct {
	Direct []*RaydiumInstructionWrap
	Inner  []*RaydiumInstructionWrap
}

func (is *RaydiumInstructions) Empty() bool {
	return len(is.Direct) == 0 && len(is.Inner) == 0
}

func GetRaydiumInstructions(tx *rpc.GetBlockTransaction) *RaydiumInstructions {
	return &RaydiumInstructions{
		Direct: parseDirectInstructions(tx),
		Inner:  parseInnerInstructions(tx),
	}
}

func parseDirectInstructions(tx *rpc.GetBlockTransaction) []*RaydiumInstructionWrap {
	direct := make([]*RaydiumInstructionWrap, 0)

	for instructionIndex, instructionAny := range tx.Transaction.Message.Instructions {
		instruction, ok := ParseInstruction(instructionAny)
		if !ok {
			continue
		}

		if filters.FilterInstructionByProgramId(instruction.ProgramId) {
			continue
		}

		instruction.Index = instructionIndex * InstructionIndexBase

		wrap := RaydiumInstructionWrap{
			Raydium:   instruction,
			Transfers: make([]*InstructionParsed, 0, 2),
		}

		for _, innerInstruction := range tx.Meta.InnerInstructions {
			if innerInstruction.Index == uint64(instructionIndex) {
				for _, innerInstructionInstruction := range innerInstruction.Instructions {
					instructionParsed, ok := ParseInstructionParsed(innerInstructionInstruction)
					if !ok {
						continue
					}
					if !instructionParsed.Check() {
						continue
					}
					wrap.Transfers = append(wrap.Transfers, instructionParsed)
				}
				break
			}
		}

		if len(wrap.Transfers) == 2 {
			direct = append(direct, &wrap)
		} else {
			//log.Logger.Warn("wrong inner instruction", zap.String("tx_hash", tx.Transaction.Signatures[0]))
		}
	}

	return direct
}

func parseInnerInstructions(tx *rpc.GetBlockTransaction) []*RaydiumInstructionWrap {
	result := make([]*RaydiumInstructionWrap, 0)

	for _, innerIx := range tx.Meta.InnerInstructions {
		var currentBase *RaydiumInstructionWrap
		transferCount := 0

		for index, ix := range innerIx.Instructions {
			// 尝试解析为基础指令
			if instruction, ok := ParseInstruction(ix); ok {
				if !filters.FilterInstructionByProgramId(instruction.ProgramId) {
					instruction.Index = int(innerIx.Index)*InstructionIndexBase + index
					// 如果之前有未完成的指令组，先清理掉
					currentBase = &RaydiumInstructionWrap{
						Raydium:   instruction,
						Transfers: make([]*InstructionParsed, 2), // 预分配2个位置
					}
					transferCount = 0
					continue
				}
			}

			// 如果有当前基础指令，尝试解析转账指令
			if currentBase != nil {
				if instruction, ok := ParseInstructionParsed(ix); ok && instruction.Check() {
					if instruction.StackHeight == currentBase.Raydium.StackHeight+1 {
						currentBase.Transfers[transferCount] = instruction
						transferCount++

						// 如果收集到两个转账指令，将当前指令组添加到结果中
						if transferCount == 2 {
							result = append(result, currentBase)
							currentBase = nil
						}
					}
				}
			}
		}
	}

	return result
}
