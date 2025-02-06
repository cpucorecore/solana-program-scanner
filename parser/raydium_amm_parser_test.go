package parser_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"solana-program-scanner/log"
	"solana-program-scanner/parser"
	"solana-program-scanner/types"
	"solana-program-scanner/types/orms"
)

func TestParseIxBuyOrSell(t *testing.T) {
	log.InitLoggerForTest()
	type Input struct {
		tx           *orms.Tx
		Amount0      string
		Amount1      string
		Destination0 string
		Market       *types.Market
	}

	tests := []struct {
		input    *Input
		expected *parser.TransferDetail
	}{
		// https://solscan.io/tx/3Qw8rKGpxUsHxU9p1WpqbxERqmePFsd3k89CyELADdM6UXSLtzyczBHqFLjHV4VZMnccq1jX8ujd3QNraxzwvov7
		{
			input: &Input{
				&orms.Tx{},
				"2475",
				"649070985389",
				"EXM9vDSud7SG5XD1AQRSdETbzLAgr7pmNjpMrXCWN6D7",
				&types.Market{
					BaseDecimal:  9,
					QuoteDecimal: 6,
					BaseVault:    "EXM9vDSud7SG5XD1AQRSdETbzLAgr7pmNjpMrXCWN6D7",
					QuoteVault:   "FNmbMQegyBptyBpixrPAp1vgkoaCMZnsvrzzhyRcm6ym",
					BaseMint:     "So11111111111111111111111111111111111111112",
					QuoteMint:    "F8T1Vrna8k3gcy9uZfNPRhRwQwZPJ47nZmMYqX8Epump",
				},
			},
			expected: &parser.TransferDetail{
				Event:         types.Buy,
				Token0Amount:  parser.MustToDecimal("649070985389", -6),
				Token0Address: "F8T1Vrna8k3gcy9uZfNPRhRwQwZPJ47nZmMYqX8Epump",
				Token1Amount:  parser.MustToDecimal("2475", -9),
				Token1Address: "So11111111111111111111111111111111111111112",
			},
		},
		// https://solscan.io/tx/5EKhSGjhPo3UXDbSBEXhazqXE1bRg6eFp9oZgA8jUgmnBM5sSeXBPqnxQCXSkTJEMvYVWx2t42TrwQqMpdX64kqE
		{
			input: &Input{
				&orms.Tx{},
				"4317688",
				"9946",
				"2CVuJtvWr4BWNHiG3zcBbseE9Jf7dmQim753aar7mRtX",
				&types.Market{
					BaseDecimal:  9,
					QuoteDecimal: 6,
					BaseVault:    "4nzsTYZeqc28yAro3euyBLdoZZK6SW2sywGpTTwvHapN",
					QuoteVault:   "2CVuJtvWr4BWNHiG3zcBbseE9Jf7dmQim753aar7mRtX",
					BaseMint:     "So11111111111111111111111111111111111111112",
					QuoteMint:    "Av67SYEwmfaPmS3Yc9gJSrGu71EQHn7PyLvEFDxtpump",
				},
			},
			expected: &parser.TransferDetail{
				Event:         types.Sell,
				Token0Amount:  parser.MustToDecimal("4317688", -6),
				Token0Address: "Av67SYEwmfaPmS3Yc9gJSrGu71EQHn7PyLvEFDxtpump",
				Token1Amount:  parser.MustToDecimal("9946", -9),
				Token1Address: "So11111111111111111111111111111111111111112",
			},
		},
		// https://solscan.io/tx/FsEoYDxhekaRxjbEr3qJfsER32sp8yNEm6adupeDnXFE6V5y2jqhLpbjHSWdxVqf7ZmxTbZhYMYUaPbnAT1K3TG
		{
			input: &Input{
				&orms.Tx{},
				"12266",
				"19033421",
				"4h9ciVLGRAiZ6G6DhkGz33BQH3VFsZSoRmhmfqWriaYS",
				&types.Market{
					BaseDecimal:  9,
					QuoteDecimal: 6,
					BaseVault:    "4h9ciVLGRAiZ6G6DhkGz33BQH3VFsZSoRmhmfqWriaYS",
					QuoteVault:   "9QCAdf7XeCzv5nsaLV3YqfRptgwJ2GpMU86bZvCE87aS",
					BaseMint:     "So11111111111111111111111111111111111111112",
					QuoteMint:    "AQfGBYk79q3KkXJtV5Jagnq59HGCb9gvnz1FZHKopump",
				},
			},
			expected: &parser.TransferDetail{
				Event:         types.Buy,
				Token0Amount:  parser.MustToDecimal("19033421", -6),
				Token0Address: "AQfGBYk79q3KkXJtV5Jagnq59HGCb9gvnz1FZHKopump",
				Token1Amount:  parser.MustToDecimal("12266", -9),
				Token1Address: "So11111111111111111111111111111111111111112",
			},
		},
		// https://solscan.io/tx/5Bce5HWPAGgZPto9nyhCxBBpuQE8apfoqsAjvdB1mScLmRQk5Towio2TMdxdJx3YaGWbnivTm6t3ooMiPwZtjimH
		{
			input: &Input{
				&orms.Tx{},
				"903538002",
				"1186160826",
				"5UWmR7fkMLbqqVAHCLfDzPaLaQSZ1scqJPEDWHEhS79s",
				&types.Market{
					BaseDecimal:  9,
					QuoteDecimal: 6,
					BaseVault:    "4U941P9qypjiVmnekg8XGrvzgzYECBd65wmAETLwDa6u",
					QuoteVault:   "5UWmR7fkMLbqqVAHCLfDzPaLaQSZ1scqJPEDWHEhS79s",
					BaseMint:     "So11111111111111111111111111111111111111112",
					QuoteMint:    "A8C3xuqscfmyLrte3VmTqrAq8kgMASius9AFNANwpump",
				},
			},
			expected: &parser.TransferDetail{
				Event:         types.Sell,
				Token0Amount:  parser.MustToDecimal("903538002", -6),
				Token0Address: "A8C3xuqscfmyLrte3VmTqrAq8kgMASius9AFNANwpump",
				Token1Amount:  parser.MustToDecimal("1186160826", -9),
				Token1Address: "So11111111111111111111111111111111111111112",
			},
		},
	}

	for tid, test := range tests {
		transferDetail := parser.ParseEventBuyOrSell(test.input.Amount0, test.input.Amount1, test.input.Destination0, test.input.Market, decimal.NewFromFloat(200.0))
		require.Equal(t, test.expected.Event, transferDetail.Event, tid)
		require.Equal(t, test.expected.Token0Address, transferDetail.Token0Address, tid)
		require.Equal(t, test.expected.Token0Amount, transferDetail.Token0Amount, tid)
		require.Equal(t, test.expected.Token1Address, transferDetail.Token1Address, tid)
		require.Equal(t, test.expected.Token1Amount, transferDetail.Token1Amount, tid)
	}
}

func TestParseEventPoolRemove(t *testing.T) {
	tests := []struct {
		name     string
		amount0  string
		amount1  string
		source0  string
		market   *types.Market
		expected *parser.TransferDetail
	}{
		{
			name:    "Remove liquidity when source0 is BaseVault and base is SOL",
			amount0: "1000000000", // 1 SOL
			amount1: "2000000000", // 2 USDC
			source0: "base_vault_1",
			market: &types.Market{
				BaseMint:     types.SOL,       // SOL is base token
				QuoteMint:    "quote_mint_1",  // USDC mint
				BaseVault:    "base_vault_1",  // SOL vault
				QuoteVault:   "quote_vault_1", // USDC vault
				BaseDecimal:  9,
				QuoteDecimal: 9,
			},
			expected: &parser.TransferDetail{
				Token0Amount:  decimal.NewFromFloat(2), // USDC amount
				Token0Address: "quote_mint_1",          // USDC mint
				Token1Amount:  decimal.NewFromFloat(1), // SOL amount
				Token1Address: types.SOL,               // SOL mint
				AmountUsd:     decimal.Zero,
				PriceUsd:      decimal.Zero,
			},
		},
		{
			name:    "Remove liquidity when source0 is QuoteVault and quote is SOL",
			amount0: "2000000000", // 2 TOKEN
			amount1: "1000000000", // 1 SOL
			source0: "quote_vault_2",
			market: &types.Market{
				BaseMint:     "base_mint_2",   // Token mint
				QuoteMint:    types.SOL,       // SOL is quote token
				BaseVault:    "base_vault_2",  // Token vault
				QuoteVault:   "quote_vault_2", // SOL vault
				BaseDecimal:  9,
				QuoteDecimal: 9,
			},
			expected: &parser.TransferDetail{
				Token0Amount:  decimal.NewFromFloat(2), // Token amount
				Token0Address: "base_mint_2",           // Token mint
				Token1Amount:  decimal.NewFromFloat(1), // SOL amount
				Token1Address: types.SOL,               // SOL mint
				AmountUsd:     decimal.Zero,
				PriceUsd:      decimal.Zero,
			},
		},
		{
			name:    "Remove liquidity when source0 is BaseVault and base is Token",
			amount0: "2000000",    // 2 USDC
			amount1: "1000000000", // 1 SOL
			source0: "base_vault_3",
			market: &types.Market{
				BaseMint:     "base_mint_3",   // USDC mint
				QuoteMint:    types.SOL,       // SOL is quote token
				BaseVault:    "base_vault_3",  // USDC vault
				QuoteVault:   "quote_vault_3", // SOL vault
				BaseDecimal:  6,
				QuoteDecimal: 9,
			},
			expected: &parser.TransferDetail{
				Token0Amount:  decimal.NewFromFloat(2), // USDC amount
				Token0Address: "base_mint_3",           // USDC mint
				Token1Amount:  decimal.NewFromFloat(1), // SOL amount
				Token1Address: types.SOL,               // SOL mint
				AmountUsd:     decimal.Zero,
				PriceUsd:      decimal.Zero,
			},
		},
		{
			name:    "Remove liquidity when source0 is QuoteVault and quote is Token",
			amount0: "1000000000", // 1 SOL
			amount1: "2000000",    // 2 USDC
			source0: "quote_vault_4",
			market: &types.Market{
				BaseMint:     types.SOL,       // SOL is base token
				QuoteMint:    "quote_mint_4",  // USDC mint
				BaseVault:    "base_vault_4",  // SOL vault
				QuoteVault:   "quote_vault_4", // USDC vault
				BaseDecimal:  9,
				QuoteDecimal: 6,
			},
			expected: &parser.TransferDetail{
				Token0Amount:  decimal.NewFromFloat(2), // USDC amount
				Token0Address: "quote_mint_4",          // USDC mint
				Token1Amount:  decimal.NewFromFloat(1), // SOL amount
				Token1Address: types.SOL,               // SOL mint
				AmountUsd:     decimal.Zero,
				PriceUsd:      decimal.Zero,
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(i)
			result := parser.ParseEventPoolRemove(tt.amount0, tt.amount1, tt.source0, tt.market)

			assert.Equal(t, tt.expected.Token0Amount.String(), result.Token0Amount.String())
			assert.Equal(t, tt.expected.Token0Address, result.Token0Address)
			assert.Equal(t, tt.expected.Token1Amount.String(), result.Token1Amount.String())
			assert.Equal(t, tt.expected.Token1Address, result.Token1Address)
			assert.Equal(t, tt.expected.AmountUsd.String(), result.AmountUsd.String())
			assert.Equal(t, tt.expected.PriceUsd.String(), result.PriceUsd.String())
		})
	}
}
