package market_detail_getter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBalanceWithRetry(t *testing.T) {
	getter := NewAccountBalanceGetter("https://api.mainnet-beta.solana.com")
	accountBalance, err := getter.GetBalanceWithRetry("5yGRKYyxuHuPQAdKrBC1C5NfanG7ZApgHvAAJ1GPtDWt")
	require.Nil(t, err)
	fmt.Println(accountBalance)
}
