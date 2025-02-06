package parser

import "encoding/json"

/*
{
    "parsed": {
        "info": {
            "amount": "1811260",
            "authority": "6NVomBM5LJmfkMhKDH6BJrXD8xCJAta4wS4L44xba22J",
            "destination": "pn6a54Jqme6X3G2dy1FkExKr4f5xtdT8nrCZJQWV7iF",
            "source": "579yCFDRu7U6pxYq8ieE7QBG3JYiHiJG3JhpDzeZfaNy"
        },
        "type": "transfer"
    },
    "program": "spl-token",
    "programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
    "stackHeight": 2
}
*/

type InstructionParsed struct {
	Parsed struct {
		Info struct {
			Amount      string `json:"amount"`
			Authority   string `json:"authority"`
			Destination string `json:"destination"`
			Source      string `json:"source"`
		} `json:"info"`
		Type string `json:"type"`
	} `json:"parsed"`
	Program     string `json:"program"`
	ProgramId   string `json:"programId"`
	StackHeight int    `json:"stackHeight"`
}

func UnmarshalInstructionParsed(instructions any) (*InstructionParsed, bool) {
	bs, err := json.Marshal(instructions)
	if err != nil {
		return nil, false
	}
	var tx InstructionParsed
	err = json.Unmarshal(bs, &tx)
	if err != nil {
		return nil, false
	}
	return &tx, true
}

func ParseInstructionParsed(instruction any) (*InstructionParsed, bool) {
	// 将interface{}转换为map
	m, ok := instruction.(map[string]interface{})
	if !ok {
		return nil, false
	}

	result := &InstructionParsed{}

	// 解析parsed字段
	parsed, ok := m["parsed"].(map[string]interface{})
	if !ok {
		return nil, false
	}

	// 解析info字段
	info, ok := parsed["info"].(map[string]interface{})
	if !ok {
		return nil, false
	}

	// 解析info中的字段
	if amount, ok := info["amount"].(string); ok {
		result.Parsed.Info.Amount = amount
	}
	if authority, ok := info["authority"].(string); ok {
		result.Parsed.Info.Authority = authority
	}
	if destination, ok := info["destination"].(string); ok {
		result.Parsed.Info.Destination = destination
	}
	if source, ok := info["source"].(string); ok {
		result.Parsed.Info.Source = source
	}
	// 解析type字段
	if typeStr, ok := parsed["type"].(string); ok {
		result.Parsed.Type = typeStr
	}

	// 解析program字段
	if program, ok := m["program"].(string); ok {
		result.Program = program
	}

	// 解析programId字段
	if programId, ok := m["programId"].(string); ok {
		result.ProgramId = programId
	}

	// 解析stackHeight字段
	if stackHeight, ok := m["stackHeight"].(float64); ok {
		result.StackHeight = int(stackHeight)
	}

	return result, true
}

func (i *InstructionParsed) Check() bool {
	if i.Program != "spl-token" {
		return false
	}
	if i.Parsed.Type != "transfer" {
		return false
	}
	return true
}
