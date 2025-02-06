package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInnerInstructionParsed(t *testing.T) {
	// Test case 1: Valid JSON input
	validJSON := `{
		"index": 2,
		"instructions": [
			{
				"parsed": {
					"info": {
						"amount": "1803557",
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
		]
	}`

	var validData interface{}
	err := json.Unmarshal([]byte(validJSON), &validData)
	assert.NoError(t, err)

	result, ok := ParseInnerInstructionParsed(validData)
	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Index)
	assert.Len(t, result.Instructions, 1)
	assert.Equal(t, "1803557", result.Instructions[0].Parsed.Info.Amount)
	assert.Equal(t, "transfer", result.Instructions[0].Parsed.Type)
	assert.Equal(t, "spl-token", result.Instructions[0].Program)

	// Test case 2: Invalid input (not a map)
	invalidData := "not a map"
	result, ok = ParseInnerInstructionParsed(invalidData)
	assert.False(t, ok)
	assert.Nil(t, result)

	// Test case 3: Missing instructions field
	invalidJSON := `{
		"index": 2
	}`
	var missingInstructionsData interface{}
	err = json.Unmarshal([]byte(invalidJSON), &missingInstructionsData)
	assert.NoError(t, err)

	result, ok = ParseInnerInstructionParsed(missingInstructionsData)
	assert.False(t, ok)
	assert.Nil(t, result)

	// Test case 4: Multiple instructions
	multipleInstructionsJSON := `{
		"index": 2,
		"instructions": [
			{
				"parsed": {
					"info": {
						"amount": "1803557",
						"authority": "6NVomBM5LJmfkMhKDH6BJrXD8xCJAta4wS4L44xba22J",
						"destination": "pn6a54Jqme6X3G2dy1FkExKr4f5xtdT8nrCZJQWV7iF",
						"source": "579yCFDRu7U6pxYq8ieE7QBG3JYiHiJG3JhpDzeZfaNy"
					},
					"type": "transfer"
				},
				"program": "spl-token",
				"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
				"stackHeight": 2
			},
			{
				"parsed": {
					"info": {
						"amount": "3075499413650",
						"authority": "5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1",
						"destination": "7UaTBNd6Hz2Aq8nz95nGLdAdS7qFvqhe84LXqf2QVa1A",
						"source": "6Nzx6JfR2d6hCPtHDTU6v1b2RAcaeNPNhigaUM771Yfe"
					},
					"type": "transfer"
				},
				"program": "spl-token",
				"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
				"stackHeight": 2
			}
		]
	}`

	var multipleInstructionsData interface{}
	err = json.Unmarshal([]byte(multipleInstructionsJSON), &multipleInstructionsData)
	assert.NoError(t, err)

	result, ok = ParseInnerInstructionParsed(multipleInstructionsData)
	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Index)
	assert.Len(t, result.Instructions, 2)
	assert.Equal(t, "1803557", result.Instructions[0].Parsed.Info.Amount)
	assert.Equal(t, "3075499413650", result.Instructions[1].Parsed.Info.Amount)
}
