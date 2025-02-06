package parser

import (
	"encoding/json"
	"github.com/mr-tron/base58"
)

type Instruction struct {
	Accounts    []string
	Data        string
	ProgramId   string
	StackHeight int

	Index     int
	DataBytes []byte
}

func UnmarshalInstruction(instructions any) (*Instruction, bool) {
	bs, err := json.Marshal(instructions)
	if err != nil {
		return nil, false
	}
	var tx Instruction
	err = json.Unmarshal(bs, &tx)
	if err != nil {
		return nil, false
	}
	return &tx, true
}

func ParseInstruction(input interface{}) (*Instruction, bool) {
	// 将interface{}转换为map
	m, ok := input.(map[string]interface{})
	if !ok {
		return nil, false
	}

	instruction := &Instruction{}

	// 解析Accounts
	if accounts, ok := m["accounts"].([]interface{}); ok {
		instruction.Accounts = make([]string, len(accounts))
		for i, acc := range accounts {
			if accStr, ok := acc.(string); ok {
				instruction.Accounts[i] = accStr
			}
		}
	}

	// 解析Data
	if data, ok := m["data"].(string); ok {
		instruction.Data = data
		dataBytes, err := base58.Decode(data)
		if err != nil {
			return nil, false
		}
		instruction.DataBytes = dataBytes
	}

	// 解析ProgramId
	if programId, ok := m["programId"].(string); ok {
		instruction.ProgramId = programId
	}

	// 解析StackHeight
	// 注意：stackHeight可能是null或者数字
	if stackHeight, ok := m["stackHeight"].(float64); ok {
		instruction.StackHeight = int(stackHeight)
	}

	return instruction, true
}
