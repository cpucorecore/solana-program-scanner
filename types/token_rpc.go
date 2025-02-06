package types

type Token struct {
	Address         string
	Name            string
	Symbol          string
	Decimals        int
	Supply          string
	UpdateAuthority string
	Uri             string
}

type TokenResp struct {
	Status  int
	ErrCode int
	ErrMsg  string
	Token   Token
}
