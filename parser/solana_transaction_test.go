package parser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionAccountBook(t *testing.T) {
	testCases := []struct {
		name         string
		filepath     string
		accountLen   int
		firstAccount string
		lastAccount  string
	}{
		{
			name:         "Add transaction",
			filepath:     "txs/2EBaHmULTBgt7jJVYeffrMZhJsMK4R7uDfQQ3L3nn2bRToH28PY1M9deFkhyFCu2LrLsFhgGqfuJfx1Jsyo1SijX.add.json",
			accountLen:   19,
			firstAccount: "7xCgwwAhR1jjNmGyVaj4VP716H8eqaUYan6fen8JuBwP",
			lastAccount:  "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		},
		{
			name:         "Remove transaction",
			filepath:     "txs/3EKnhqM5cnLZEDAK3tMD4gTtZxD6LGqehEYQ6b4LZFCWM5Vm41WBgQBLSje5XEokFwj6iJmcZMzU7UB6ZxTrZxUc.remove.json",
			accountLen:   27,
			firstAccount: "8yvHN3WMyVWe7BNFscUAJhwdW5M3AfdY7if7PofTUJjg",
			lastAccount:  "2khkbtQSxeNGYAitynFhHdMEbm5zb87emuaDzjSDKcMB",
		},
		{
			name:         "Swap transaction",
			filepath:     "txs/55XSjkix3RbZXbRsvj8gZqdBNrUomfoiRWKpYAaRXymgSf8eyAnv7VqmTDVpRzbAPAKEX3eicGpPJh7kYhH9AHV2.swap.json",
			accountLen:   24,
			firstAccount: "7ssAF9F74txZeY7V9dyorGRusAaSyd9csYYBDxVCHDxs",
			lastAccount:  "srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX",
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
			accountBook := CreateTransactionAccountBook(&response.Result)

			// Test that we got non-empty results
			assert.NotEmpty(t, accountBook)

			// Test specific expectations
			assert.Equal(t, tc.accountLen, len(accountBook), "Account length mismatch")
			assert.Contains(t, accountBook, tc.firstAccount, "First account not found")
			assert.Contains(t, accountBook, tc.lastAccount, "Last account not found")

			// Test that all accounts in accountBook are also in accountBookWithIndex
			for addr := range accountBook {
				accountInfo, exists := accountBook[addr]
				assert.True(t, exists)
				assert.NotNil(t, accountInfo)
				assert.GreaterOrEqual(t, accountInfo.Index, 0)
			}
		})
	}
}
