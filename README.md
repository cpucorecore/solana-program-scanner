# sol-usd price pair
## query price-feed-id
[ref](https://www.pyth.network/developers/price-feed-ids)

## query price
[ref](https://docs.pyth.network/price-feeds/fetch-price-updates#rest-api)

## query sol-usd pair
```bash
curl -X 'GET' 'https://hermes.pyth.network/v2/updates/price/latest?ids%5B%5D=0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d' | jq
```
