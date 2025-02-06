package parser

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalInstruction(t *testing.T) {
	// 测试输入JSON
	jsonStr := `{
        "accounts": [
            "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
            "4a1giKd24E3t86hyyJ7NeUwjhPZFXrZFyTWrmSVbdd8Z",
            "5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1"
        ],
        "data": "5yGUcLihb5D6LmzVFph8gET",
        "programId": "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8",
        "stackHeight": null
    }`

	var input interface{}
	err := json.Unmarshal([]byte(jsonStr), &input)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	instruction, ok := UnmarshalInstruction(input)
	if !ok {
		t.Fatal("Failed to parse instruction")
	}

	// 验证解析结果
	if instruction.Data != "5yGUcLihb5D6LmzVFph8gET" {
		t.Errorf("Expected data '5yGUcLihb5D6LmzVFph8gET', got '%s'", instruction.Data)
	}

	if instruction.ProgramId != "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8" {
		t.Errorf("Expected programId '675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8', got '%s'", instruction.ProgramId)
	}

	if len(instruction.Accounts) != 3 {
		t.Errorf("Expected 3 accounts, got %d", len(instruction.Accounts))
	}

	// stackHeight应该是0，因为输入是null
	if instruction.StackHeight != 0 {
		t.Errorf("Expected stackHeight 0, got %d", instruction.StackHeight)
	}
}
