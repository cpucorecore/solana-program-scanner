package parser

type InnerInstructionParsed struct {
	Index        int
	Instructions []InstructionParsed
}

func ParseInnerInstructionParsed(v any) (*InnerInstructionParsed, bool) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, false
	}

	result := &InnerInstructionParsed{}

	// Parse index field
	if index, ok := m["index"].(float64); ok {
		result.Index = int(index)
	}

	// Parse instructions array
	instructions, ok := m["instructions"].([]interface{})
	if !ok {
		return nil, false
	}

	// Parse each instruction in the array
	result.Instructions = make([]InstructionParsed, 0, len(instructions))
	for _, inst := range instructions {
		if parsed, ok := ParseInstructionParsed(inst); ok {
			result.Instructions = append(result.Instructions, *parsed)
		}
	}

	return result, true
}
