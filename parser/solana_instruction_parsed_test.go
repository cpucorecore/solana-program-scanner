package parser

import (
	"encoding/json"
	"testing"
)

func TestParseInstructionParsed(t *testing.T) {
	// 测试输入JSON
	jsonStr := `{
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
    }`

	var input interface{}
	err := json.Unmarshal([]byte(jsonStr), &input)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	result, ok := ParseInstructionParsed(input)
	if !ok {
		t.Fatal("Failed to parse instruction")
	}
	// 验证解析结果
	if result.Parsed.Info.Amount != "1811260" {
		t.Errorf("Expected amount '1811260', got '%s'", result.Parsed.Info.Amount)
	}
	if result.Parsed.Info.Authority != "6NVomBM5LJmfkMhKDH6BJrXD8xCJAta4wS4L44xba22J" {
		t.Errorf("Expected authority '6NVomBM5LJmfkMhKDH6BJrXD8xCJAta4wS4L44xba22J', got '%s'", result.Parsed.Info.Authority)
	}
	if result.Parsed.Type != "transfer" {
		t.Errorf("Expected type 'transfer', got '%s'", result.Parsed.Type)
	}
	if result.Program != "spl-token" {
		t.Errorf("Expected program 'spl-token', got '%s'", result.Program)
	}
	if result.ProgramId != "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA" {
		t.Errorf("Expected programId 'TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA', got '%s'", result.ProgramId)
	}
	if result.StackHeight != 2 {
		t.Errorf("Expected stackHeight 2, got %d", result.StackHeight)
	}
}
