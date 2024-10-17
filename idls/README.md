```bash
git clone -b "v0.3.1" https://github.com/gagliardetto/anchor-go.git
cd cnchor-go
go build

./anchor-go --type-id uint8 --src=<the idl json file> #for raydium_amm
./anchor-go --type-id anchor --src=<the idl json file> #for openbook_v2
```

# idl source
program | url
-- | --
openbook_v2 | https://github.com/openbook-dex/openbook-v2/blob/master/idl/openbook_v2.json
raydium_amm | https://github.com/raydium-io/raydium-idl/blob/master/raydium_amm/idl.json
