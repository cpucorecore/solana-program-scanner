package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
	"solana-program-scanner/types"
)

func TestFilterIxByMarket(t *testing.T) {
	m1 := &types.Market{
		BaseMint:  types.SOL,
		QuoteMint: types.SOL,
	}
	m2 := &types.Market{
		BaseMint:  types.SOL,
		QuoteMint: "mockMint",
	}
	m3 := &types.Market{
		BaseMint:  "mockMint",
		QuoteMint: types.SOL,
	}
	m4 := &types.Market{
		BaseMint:  "mockMint",
		QuoteMint: "mockMint",
	}
	require.Equal(t, true, FilterIxByMarket(m1))
	require.Equal(t, false, FilterIxByMarket(m2))
	require.Equal(t, false, FilterIxByMarket(m3))
	require.Equal(t, true, FilterIxByMarket(m4))
}
