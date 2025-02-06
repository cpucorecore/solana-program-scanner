package parser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/stretchr/testify/assert"
)

func TestParseInnerInstructions(t *testing.T) {
	// Load test data
	jsonFile, err := os.ReadFile("tx_samples/inner/4ADotFjPe3czrm17g7sdh78aS4s5Gvz2ZcyYWUqhxeQnNstozRQwxNVQom5aKU99CRe2WBx1M48AnTbvdTikEf9C.json")
	assert.NoError(t, err)

	var tx rpc.GetBlockTransaction
	err = json.Unmarshal(jsonFile, &tx)
	assert.NoError(t, err)

	// Execute
	result := parseInnerInstructions(&tx)

	// Assert
	assert.Len(t, result, 4, "should have 4 Raydium instructions")

	// Test first instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[0].Raydium.ProgramId)
	assert.Equal(t, "64hrY8hBuRh5YBRdr9kEoxT", result[0].Raydium.Data)
	assert.Len(t, result[0].Transfers, 2)
	assert.Equal(t, "65567050", result[0].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "43124584718043", result[0].Transfers[1].Parsed.Info.Amount)

	// Test second instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[1].Raydium.ProgramId)
	assert.Equal(t, "6NdPA6aA9jM33qSHNhxJHQs", result[1].Raydium.Data)
	assert.Len(t, result[1].Transfers, 2)
	assert.Equal(t, "43124584718043", result[1].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "135686765263", result[1].Transfers[1].Parsed.Info.Amount)

	// Test third instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[2].Raydium.ProgramId)
	assert.Equal(t, "6MADnqxMrdDRYqVocpTpjgj", result[2].Raydium.Data)
	assert.Len(t, result[2].Transfers, 2)
	assert.Equal(t, "135686765263", result[2].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "122856362", result[2].Transfers[1].Parsed.Info.Amount)

	// Test fourth instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[3].Raydium.ProgramId)
	assert.Equal(t, "6GaaEuTZs6ZMuEKsnDTrUmu", result[3].Raydium.Data)
	assert.Len(t, result[3].Transfers, 2)
	assert.Equal(t, "122856362", result[3].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "65967448", result[3].Transfers[1].Parsed.Info.Amount)

	// Test transfer details for first instruction
	firstTransfer := result[0].Transfers[0]
	assert.Equal(t, "B3aD5mvj137xAtijz1nStLGdjVFLxKrGHjgg8xNTrA8R", firstTransfer.Parsed.Info.Source)
	assert.Equal(t, "Emo5LenPc2YQiABVMwmkWui6WAyBvNdN5CsxjuT6DyDy", firstTransfer.Parsed.Info.Destination)
	assert.Equal(t, "MasKodVPzTeFikftwUmW7VwJwgYMZhsATDn8H23h9RL", firstTransfer.Parsed.Info.Authority)
}

func TestParseInnerInstructions2(t *testing.T) {
	// Load test data
	jsonData, err := os.ReadFile("tx_samples/inner/3aPjKXbcr5BFKzsaUvG3vaUNVoi3rBb366JxVX4W9qCwBkLdGvTn5pqJZ6CeSqMMateQXxNyCkRYmyB9R83PiE7Q.json")
	assert.NoError(t, err)

	var tx rpc.GetBlockTransaction
	err = json.Unmarshal(jsonData, &tx)
	assert.NoError(t, err)

	// Execute
	result := parseInnerInstructions(&tx)

	// Assert
	assert.Len(t, result, 2, "should have 2 Raydium instructions")

	// Test first instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[0].Raydium.ProgramId)
	assert.Equal(t, "65QMeEFAH6ptZcQKwegT5nX", result[0].Raydium.Data)
	assert.Len(t, result[0].Transfers, 2)
	assert.Equal(t, "9970000", result[0].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "38600584", result[0].Transfers[1].Parsed.Info.Amount)

	// Test second instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[1].Raydium.ProgramId)
	assert.Equal(t, "6CRdrKDVe8DqRUJBVWk2kP1", result[1].Raydium.Data)
	assert.Len(t, result[1].Transfers, 2)
	assert.Equal(t, "38600584", result[1].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "3539168265", result[1].Transfers[1].Parsed.Info.Amount)
}

// ... existing code ...

func TestDirectInstructions(t *testing.T) {
	// Load test data
	jsonData, err := os.ReadFile("tx_samples/direct/4SQEp3QbTgQEXYrc7V1DzTVLSMxWus126JnzKnpwXrKAVhUbWuYM6YfeCcULHSNvdL4QgmnuZCX1jRmLemWuknhE.json")
	assert.NoError(t, err)

	var tx rpc.GetBlockTransaction
	err = json.Unmarshal(jsonData, &tx)
	assert.NoError(t, err)

	// Execute
	result := parseDirectInstructions(&tx)

	// Assert
	assert.Len(t, result, 6, "should have 6 Raydium instructions")

	// Test first instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[0].Raydium.ProgramId)
	assert.Equal(t, "6Su5z74bs4iFRtjUNno97Ew", result[0].Raydium.Data)
	assert.Len(t, result[0].Transfers, 2)
	assert.Equal(t, "7156990", result[0].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "483", result[0].Transfers[1].Parsed.Info.Amount)

	// Test second instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[1].Raydium.ProgramId)
	assert.Equal(t, "7NxfptFcrZdX5Qz8KeuudVR", result[1].Raydium.Data)
	assert.Len(t, result[1].Transfers, 2)
	assert.Equal(t, "487", result[1].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "7156990", result[1].Transfers[1].Parsed.Info.Amount)

	// Test transfer details for first instruction
	firstTransfer := result[0].Transfers[0]
	assert.Equal(t, "8sr8pNgvrvUYWkiZZfHq11nbfJFAXriB3DAij33PLykJ", firstTransfer.Parsed.Info.Source)
	assert.Equal(t, "Go8SGmgAoHFFBqYRwLoprADcJMNi3jyt2zwSnQfEdg4W", firstTransfer.Parsed.Info.Destination)
	assert.Equal(t, "EfqfyyPwQmE5foGbCopa3WrfsaHR5Ys7LocnVRd7SEsN", firstTransfer.Parsed.Info.Authority)

	// Test last instruction
	lastIdx := len(result) - 1
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[lastIdx].Raydium.ProgramId)
	assert.Equal(t, "7NxfptFcrZdX5Qz8KeuudVR", result[lastIdx].Raydium.Data)
	assert.Len(t, result[lastIdx].Transfers, 2)
	assert.Equal(t, "487", result[lastIdx].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "7156990", result[lastIdx].Transfers[1].Parsed.Info.Amount)
}

func TestDirectInstructions2(t *testing.T) {
	// Load test data
	jsonData, err := os.ReadFile("tx_samples/direct/4sqhEHsoHvPXmgUiTBxKf6iwnsX8ccqesavBpyQisQJy8WfL6dgNLVuae6Z3ppCZLDA63VNXMkJou3VG6tpF9QPL.json")
	assert.NoError(t, err)

	var tx rpc.GetBlockTransaction
	err = json.Unmarshal(jsonData, &tx)
	assert.NoError(t, err)

	// Execute
	result := parseDirectInstructions(&tx)

	// Assert
	assert.Len(t, result, 1, "should have 1 Raydium instruction")

	// Test the instruction
	assert.Equal(t, "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", result[0].Raydium.ProgramId)
	assert.Equal(t, "6BQCg45VY938k1hJgBfDtzj", result[0].Raydium.Data)
	assert.Len(t, result[0].Transfers, 2)
	assert.Equal(t, "250000000", result[0].Transfers[0].Parsed.Info.Amount)
	assert.Equal(t, "1073587018694", result[0].Transfers[1].Parsed.Info.Amount)

	// Test transfer details
	firstTransfer := result[0].Transfers[0]
	assert.Equal(t, "8fFv27FrTwrGkEyhGR7QSXwGuvJNkjJxjBaaceDj6qNt", firstTransfer.Parsed.Info.Source)
	assert.Equal(t, "CDQim93Fo1RcLrCSMHuWBLZu27ZmeJnogy1G5jAa6s8Q", firstTransfer.Parsed.Info.Destination)
	assert.Equal(t, "HqSVn1cLYqVXn41cCoxwPu2unFZzk2toNgsy12jDyUV1", firstTransfer.Parsed.Info.Authority)
}
