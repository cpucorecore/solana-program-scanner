package types

type Token struct {
	Id              string
	Address         string
	Name            string
	Decimal         int
	Creator         string
	ChainId         int
	TotalSupply     string
	Symbol          string
	Block           uint64
	BlockAt         string `json:"block_at,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	Logo            string
	Website         string
	Twitter         string
	Telegram        string
	HolderCount     int64
	SortId          int64
	MarketAddresses []string
}
