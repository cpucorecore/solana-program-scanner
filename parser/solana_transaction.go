package parser

import "github.com/blocto/solana-go-sdk/rpc"

type TxAccount struct {
	Index      int
	AccountKey rpc.AccountKey
}

func CreateTransactionAccountBook(tx *rpc.GetBlockTransaction) map[string]TxAccount {
	accountBookWithIndex := make(map[string]TxAccount, len(tx.Transaction.Message.AccountKeys))
	for index, accountKey := range tx.Transaction.Message.AccountKeys {
		accountBookWithIndex[accountKey.Pubkey] = TxAccount{
			Index:      index,
			AccountKey: accountKey,
		}
	}
	return accountBookWithIndex
}
