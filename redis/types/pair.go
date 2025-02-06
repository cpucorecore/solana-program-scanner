package types

type Amount struct {
	Slot  uint64
	Value string
}

type Pair struct {
	Id        string
	Address   string
	Name      string
	Token0    string
	Token1    string
	ChainId   int
	Reserve0  *Amount
	Reserve1  *Amount
	Block     uint64
	BlockAt   string `json:"block_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	SortId    int64
}
