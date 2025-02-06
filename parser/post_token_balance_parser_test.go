package parser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountIndex(t *testing.T) {
	testCases := []struct {
		name       string
		filepath   string
		pool0Addr  string
		pool1Addr  string
		pool0Index uint64
		pool1Index uint64
	}{
		{
			name:       "Add transaction",
			filepath:   "txs/2EBaHmULTBgt7jJVYeffrMZhJsMK4R7uDfQQ3L3nn2bRToH28PY1M9deFkhyFCu2LrLsFhgGqfuJfx1Jsyo1SijX.add.json",
			pool0Addr:  "46f9VtEs54rjSMSVweDA5bKPfHXFW4go6YCVWfKj6NMd",
			pool1Addr:  "NoiaivFzNbcJ1WYCRVH1ZHAUEbJbqjx6S25vKmngcc4",
			pool0Index: 3,
			pool1Index: 8,
		},
		{
			name:       "Remove transaction",
			filepath:   "txs/3EKnhqM5cnLZEDAK3tMD4gTtZxD6LGqehEYQ6b4LZFCWM5Vm41WBgQBLSje5XEokFwj6iJmcZMzU7UB6ZxTrZxUc.remove.json",
			pool0Addr:  "jyLGCwC5Rcgu9jBMqimFhByG4p2Qes2Yu2k14tFYTan",
			pool1Addr:  "65LKUq1cjnua6D22azrBMweLTvVBuFuvbw5KKXb5n1HQ",
			pool0Index: 6,
			pool1Index: 7,
		},
		{
			name:       "Swap transaction",
			filepath:   "txs/55XSjkix3RbZXbRsvj8gZqdBNrUomfoiRWKpYAaRXymgSf8eyAnv7VqmTDVpRzbAPAKEX3eicGpPJh7kYhH9AHV2.swap.json",
			pool0Addr:  "8kJsfhGc1h899PE8bUSrdwpED48WMZqqoc12JZTzPBGV",
			pool1Addr:  "7qZn9RijS6J3br8rwmVoia7XaTmKAYPStVHZvksseQBd",
			pool0Index: 13,
			pool1Index: 14,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Load test data from file
			jsonData, err := os.ReadFile(tc.filepath)
			assert.NoError(t, err)

			// First unmarshal the full JSON RPC response
			var response struct {
				JsonRPC string                  `json:"jsonrpc"`
				Result  rpc.GetBlockTransaction `json:"result"`
				ID      int                     `json:"id"`
			}
			err = json.Unmarshal(jsonData, &response)
			assert.NoError(t, err)

			// Use the result field as the transaction
			accountBookWithIndex := CreateTransactionAccountBook(&response.Result)
			index0, ok := getAccountIndex(accountBookWithIndex, tc.pool0Addr)
			assert.True(t, ok)
			assert.Equal(t, tc.pool0Index, index0)
			index1, ok := getAccountIndex(accountBookWithIndex, tc.pool1Addr)
			assert.True(t, ok)
			assert.Equal(t, tc.pool1Index, index1)

			// Test that we got non-empty results
			assert.NotEmpty(t, accountBookWithIndex)
		})
	}
}

func TestGetBalance(t *testing.T) {
	testCases := []struct {
		name         string
		filepath     string
		pool0Addr    string
		pool1Addr    string
		pool0Balance string
		pool1Balance string
	}{
		{
			name:         "Add transaction",
			filepath:     "txs/2EBaHmULTBgt7jJVYeffrMZhJsMK4R7uDfQQ3L3nn2bRToH28PY1M9deFkhyFCu2LrLsFhgGqfuJfx1Jsyo1SijX.add.json",
			pool0Addr:    "46f9VtEs54rjSMSVweDA5bKPfHXFW4go6YCVWfKj6NMd",
			pool1Addr:    "NoiaivFzNbcJ1WYCRVH1ZHAUEbJbqjx6S25vKmngcc4",
			pool0Balance: "8510550.580528955",
			pool1Balance: "0.567833464",
		},
		{
			name:         "Remove transaction",
			filepath:     "txs/3EKnhqM5cnLZEDAK3tMD4gTtZxD6LGqehEYQ6b4LZFCWM5Vm41WBgQBLSje5XEokFwj6iJmcZMzU7UB6ZxTrZxUc.remove.json",
			pool0Addr:    "jyLGCwC5Rcgu9jBMqimFhByG4p2Qes2Yu2k14tFYTan",
			pool1Addr:    "65LKUq1cjnua6D22azrBMweLTvVBuFuvbw5KKXb5n1HQ",
			pool0Balance: "746673.297924",
			pool1Balance: "0.037164495",
		},
		{
			name:         "Swap transaction",
			filepath:     "txs/55XSjkix3RbZXbRsvj8gZqdBNrUomfoiRWKpYAaRXymgSf8eyAnv7VqmTDVpRzbAPAKEX3eicGpPJh7kYhH9AHV2.swap.json",
			pool0Addr:    "8kJsfhGc1h899PE8bUSrdwpED48WMZqqoc12JZTzPBGV",
			pool1Addr:    "7qZn9RijS6J3br8rwmVoia7XaTmKAYPStVHZvksseQBd",
			pool0Balance: "20.559584048",
			pool1Balance: "946075371.540603",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Load test data from file
			jsonData, err := os.ReadFile(tc.filepath)
			assert.NoError(t, err)

			// First unmarshal the full JSON RPC response
			var response struct {
				JsonRPC string                  `json:"jsonrpc"`
				Result  rpc.GetBlockTransaction `json:"result"`
				ID      int                     `json:"id"`
			}
			err = json.Unmarshal(jsonData, &response)
			assert.NoError(t, err)

			// Use the result field as the transaction
			accountBookWithIndex := CreateTransactionAccountBook(&response.Result)
			index0, _ := getAccountIndex(accountBookWithIndex, tc.pool0Addr)
			index1, _ := getAccountIndex(accountBookWithIndex, tc.pool1Addr)
			balance0, ok := getBalance(index0, &response.Result)
			assert.True(t, ok)
			assert.Equal(t, tc.pool0Balance, balance0)

			balance1, ok := getBalance(index1, &response.Result)
			assert.True(t, ok)
			assert.Equal(t, tc.pool1Balance, balance1)
		})
	}
}

func TestGetPostTokenBalance(t *testing.T) {
	testCases := []struct {
		name         string
		filepath     string
		pool0Addr    string
		pool1Addr    string
		pool0Balance string
		pool1Balance string
	}{
		{
			name:         "Add transaction",
			filepath:     "txs/2EBaHmULTBgt7jJVYeffrMZhJsMK4R7uDfQQ3L3nn2bRToH28PY1M9deFkhyFCu2LrLsFhgGqfuJfx1Jsyo1SijX.add.json",
			pool0Addr:    "46f9VtEs54rjSMSVweDA5bKPfHXFW4go6YCVWfKj6NMd",
			pool1Addr:    "NoiaivFzNbcJ1WYCRVH1ZHAUEbJbqjx6S25vKmngcc4",
			pool0Balance: "8510550.580528955",
			pool1Balance: "0.567833464",
		},
		{
			name:         "Remove transaction",
			filepath:     "txs/3EKnhqM5cnLZEDAK3tMD4gTtZxD6LGqehEYQ6b4LZFCWM5Vm41WBgQBLSje5XEokFwj6iJmcZMzU7UB6ZxTrZxUc.remove.json",
			pool0Addr:    "jyLGCwC5Rcgu9jBMqimFhByG4p2Qes2Yu2k14tFYTan",
			pool1Addr:    "65LKUq1cjnua6D22azrBMweLTvVBuFuvbw5KKXb5n1HQ",
			pool0Balance: "746673.297924",
			pool1Balance: "0.037164495",
		},
		{
			name:         "Swap transaction",
			filepath:     "txs/55XSjkix3RbZXbRsvj8gZqdBNrUomfoiRWKpYAaRXymgSf8eyAnv7VqmTDVpRzbAPAKEX3eicGpPJh7kYhH9AHV2.swap.json",
			pool0Addr:    "8kJsfhGc1h899PE8bUSrdwpED48WMZqqoc12JZTzPBGV",
			pool1Addr:    "7qZn9RijS6J3br8rwmVoia7XaTmKAYPStVHZvksseQBd",
			pool0Balance: "20.559584048",
			pool1Balance: "946075371.540603",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Load test data from file
			jsonData, err := os.ReadFile(tc.filepath)
			assert.NoError(t, err)

			// First unmarshal the full JSON RPC response
			var response struct {
				JsonRPC string                  `json:"jsonrpc"`
				Result  rpc.GetBlockTransaction `json:"result"`
				ID      int                     `json:"id"`
			}
			err = json.Unmarshal(jsonData, &response)
			assert.NoError(t, err)

			// Use the result field as the transaction
			accountBookWithIndex := CreateTransactionAccountBook(&response.Result)
			index0, _ := getAccountIndex(accountBookWithIndex, tc.pool0Addr)
			index1, _ := getAccountIndex(accountBookWithIndex, tc.pool1Addr)
			balance0, ok := getBalance(index0, &response.Result)
			assert.True(t, ok)
			assert.Equal(t, tc.pool0Balance, balance0)

			balance1, ok := getBalance(index1, &response.Result)
			assert.True(t, ok)
			assert.Equal(t, tc.pool1Balance, balance1)
		})
	}
}
