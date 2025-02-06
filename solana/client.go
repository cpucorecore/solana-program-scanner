package solana

import (
	"github.com/blocto/solana-go-sdk/rpc"
	"solana-program-scanner/http_client"
)

func NewGetBlockClient(endpoint string) *rpc.RpcClient {
	rpcClient := rpc.New(
		rpc.WithEndpoint(endpoint),
		rpc.WithHTTPClient(http_client.GetBlockHttpClient))
	return &rpcClient
}

func NewCommonClient(endpoint string) *rpc.RpcClient {
	rpcClient := rpc.New(
		rpc.WithEndpoint(endpoint),
		rpc.WithHTTPClient(http_client.CommonHttpClient),
	)
	return &rpcClient
}
