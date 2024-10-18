package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/mr-tron/base58"
	"time"
)

const (
	MarketAccountDataLength = 752
	BaseDecimalStartIndex   = 32
	QuoteDecimalStartIndex  = 40
	BaseMintStartIndex      = 400
	QuoteMintStartIndex     = 432
)

type MarketGetter interface {
	getMarket(marketAddress string) (market *Market, err error)
}

func (bg *BlockGetter) getMarket(marketAddress string) (market *Market, err error) {
	config := rpc.GetAccountInfoConfig{
		Encoding: rpc.AccountEncodingBase64,
		DataSlice: &rpc.DataSlice{
			Offset: 0,
			Length: MarketAccountDataLength * 3,
		},
	}

	var data any
	for {
		resp, err := bg.cli.GetAccountInfoWithConfig(context.Background(), marketAddress, config)

		if err != nil {
			bg.fc.onErr()
			continue
		}

		if resp.Error != nil {
			bg.fc.onDone(time.Now())
			break
		}

		bg.fc.onDone(time.Now())

		data = resp.Result.Value.Data
	}

	dataArray, ok := data.([]any)
	if !ok {
		Logger.Fatal(fmt.Sprintf("type assertion err:%v on 'Data'", err))
	}

	accountDataStr, ok := dataArray[0].(string) // dataArray[1] = "base64" | ...
	if !ok {
		Logger.Fatal(fmt.Sprintf("type assertion err:%v on 'Data[0]'", err))
	}

	decodeString, err := base64.StdEncoding.DecodeString(accountDataStr)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("base64 decode accountDataStr err:%v", err))
	}

	if len(decodeString) != MarketAccountDataLength {
		Logger.Fatal("wrong data length for market account")
	}

	baseDecimalBytes := decodeString[BaseDecimalStartIndex : BaseDecimalStartIndex+8]
	quoteDecimalBytes := decodeString[QuoteDecimalStartIndex : QuoteDecimalStartIndex+8]
	baseMintBytes := decodeString[BaseMintStartIndex : BaseMintStartIndex+32]
	quoteMintBytes := decodeString[QuoteMintStartIndex : QuoteMintStartIndex+32]

	return &Market{
		Address:      marketAddress,
		BaseDecimal:  getUint64ByBytesLE(baseDecimalBytes),
		QuoteDecimal: getUint64ByBytesLE(quoteDecimalBytes),
		BaseMint:     base58.Encode(baseMintBytes),
		QuoteMint:    base58.Encode(quoteMintBytes),
	}, nil
}

func getUint64ByBytesLE(bytesLE []byte) uint64 {
	var num uint64
	for i := len(bytesLE) - 1; i >= 0; i-- {
		num = num<<8 + uint64(bytesLE[i])
	}

	return num
}
