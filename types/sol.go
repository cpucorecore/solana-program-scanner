package types

const (
	SOL = "So11111111111111111111111111111111111111112"
)

func IsTokenSol(tokenId string) bool {
	return tokenId == SOL
}
