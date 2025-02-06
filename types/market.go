package types

type Market struct {
	ChainId      int    `json:"chain_id"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	BaseDecimal  int32  `json:"base_decimal"`
	QuoteDecimal int32  `json:"quote_decimal"`
	BaseVault    string `json:"base_vault"`
	QuoteVault   string `json:"quote_vault"`
	BaseMint     string `json:"base_mint"`
	QuoteMint    string `json:"quote_mint"`
}
