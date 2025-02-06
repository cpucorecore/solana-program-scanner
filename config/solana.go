package config

type SolanaConf struct {
	ChainId         int
	RPCEndpointHTTP string
	RPCEndpointWS   string
	Commitment      string
}

var defaultSolanaConf = &SolanaConf{
	ChainId:         900,
	RPCEndpointHTTP: "https://api.mainnet-beta.solana.com",
	RPCEndpointWS:   "wss://api.mainnet-beta.solana.com",
	Commitment:      "confirmed",
}
