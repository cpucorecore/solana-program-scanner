package parser

import "github.com/blocto/solana-go-sdk/rpc"

func getAccountIndex(accountBookWithIndex map[string]TxAccount, addr string) (uint64, bool) {
	for _, account := range accountBookWithIndex {
		if account.AccountKey.Pubkey == addr {
			return uint64(account.Index), true
		}
	}
	return 0, false
}

func getBalance(accountIndex uint64, tx *rpc.GetBlockTransaction) (string, bool) {
	for _, PostTokenBalance := range tx.Meta.PostTokenBalances {
		if PostTokenBalance.AccountIndex == accountIndex {
			return PostTokenBalance.UITokenAmount.UIAmountString, true
		}
	}
	return "", false
}

type PoolBalance struct {
	TxHash string
	Market string
	B0     string
	B1     string
}

func getTxPoolBalances(tx *rpc.GetBlockTransaction, accountBookWithIndex map[string]TxAccount, ixs []*RaydiumInstructionParsed) []*PoolBalance {
	poolInfo := make([]*PoolBalance, 0, len(ixs))
	for _, ix := range ixs {
		Token0VaultAccountIndex, _ := getAccountIndex(accountBookWithIndex, ix.TransferDetail.Token0Vault)
		Token1VaultAccountIndex, _ := getAccountIndex(accountBookWithIndex, ix.TransferDetail.Token1Vault)

		token0Balance, _ := getBalance(Token0VaultAccountIndex, tx)
		token1Balance, _ := getBalance(Token1VaultAccountIndex, tx)
		poolInfo = append(poolInfo, &PoolBalance{
			TxHash: tx.Transaction.Signatures[0],
			Market: ix.MarketAddress,
			B0:     token0Balance,
			B1:     token1Balance,
		})
	}

	return poolInfo
}

func mergePoolBalances(poolBalance []*PoolBalance) []*PoolBalance {
	mapMerger := make(map[string]*PoolBalance, len(poolBalance))
	for _, pool := range poolBalance {
		mapMerger[pool.Market] = pool
	}

	output := make([]*PoolBalance, 0, len(poolBalance))
	for _, pool := range mapMerger {
		output = append(output, pool)
	}
	return output
}
