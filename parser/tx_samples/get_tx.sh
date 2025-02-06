#!/bin/bash

tx_hash=$1  # 获取第一个命令行参数作为交易哈希

curl https://api.mainnet-beta.solana.com -X POST -H "Content-Type: application/json" -d "{
  \"jsonrpc\": \"2.0\",
  \"id\": 1,
  \"method\": \"getTransaction\",
  \"params\": [
    \"$tx_hash\",
    {
      \"encoding\": \"jsonParsed\",
      \"maxSupportedTransactionVersion\": 0,
      \"commitment\": \"confirmed\"
    }
  ]
}" | jq '.result' > "${tx_hash}.json"