package market_getter

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"testing"
)

func getAmmTokens(ammPubkey string) {
	// 连接到 Solana RPC
	client := rpc.New("https://api.mainnet-beta.solana.com")

	// 解析 AMM pubkey
	ammAccount, _ := solana.PublicKeyFromBase58(ammPubkey)

	// 获取 AMM 账户信息
	accountInfo, err := client.GetAccountInfo(
		context.Background(),
		ammAccount,
	)
	if err != nil {
		panic(err)
	}

	// 获取 vault 账户的 pubkey
	data := accountInfo.Value.Data.GetBinary()
	tokenAVault := solana.PublicKeyFromBytes(data[336 : 336+32]) // 具体偏移量需要根据实际结构确认
	tokenBVault := solana.PublicKeyFromBytes(data[368 : 368+32])
	//tokenAVault, _ = solana.PublicKeyFromBase58("5yGRKYyxuHuPQAdKrBC1C5NfanG7ZApgHvAAJ1GPtDWt")
	//tokenBVault, _ = solana.PublicKeyFromBase58("ANGF1opUebqn6DDHvGdb7GC2MBsrkqQaFz8e3xczNy8E")

	// 获取 vault 账户余额
	tokenAInfo, err := client.GetTokenAccountBalance(
		context.Background(),
		tokenAVault,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}
	tokenBInfo, err := client.GetTokenAccountBalance(
		context.Background(),
		tokenBVault,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Token A amount: %s\n", tokenAInfo.Value.Amount)
	fmt.Printf("Token B amount: %s\n", tokenBInfo.Value.Amount)
}

func Test_getAmmTokens(t *testing.T) {
	getAmmTokens("DPswbfa2E7B7XcyvH18YZJTaWPuyrreWyMb938goXxy8")
}
