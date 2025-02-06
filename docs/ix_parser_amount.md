获取交易所在市场的market信息

如[交易](https://solscan.io/tx/37t9cpqQfUS3PWtJcucWsWgiQin1DqQCd2rHHSLikDKG1HdYmxPDiP6V1WqUPpUftTGJkzcyBheDPxd5Ws7mXQCd)

对应的[market](https://solscan.io/account/4ei6Y8WoySLarp4rjdegKyJpGNqxs7RiYgyJck7SeRYW#data)
只需要关注这6个字段,其他忽略
```json
{

  "baseDecimal": "6",
  "quoteDecimal": "9",
  "baseVault": "FmJe8E2LssvAU4dSPPoS6iiZTySy7dae114nmDvBTh64",
  "quoteVault": "35HseUkS5TMjHeS56NTT6h9PBn1eAJ7STv5fJ2bpit9V",
  "baseMint": "AEGSZ3xkSiemPb1YzevAqQvoLSKzyCFijV7SXJv5pump",
  "quoteMint": "So11111111111111111111111111111111111111112"
}
```
- baseDecimal是token的精度,这里表示代币AEGSZ3xkSiemPb1YzevAqQvoLSKzyCFijV7SXJv5pump的精度是6
- quoteDecimal这里表示So11111111111111111111111111111111111111112的精度是9
- baseVault是base代币的池子地址,记做pool1
- quoteVault是quote代币的池子地址,记做pool2

回头看交易

看交易的第二个指令的`Instruaction Data`部分:
```json
{
  "discriminator": {
    "type": "u8",
    "data": 9
  },
  "amountIn": {
    "type": "u64",
    "data": 1000
  },
  "minimumAmountOut": {
    "type": "u64",
    "data": 0
  }
}
```

discriminator=9表示是一个swapBaseIn指令,它有2个参数:
- amountIn: 表示向池子注入的代币数量
- minimumAmountOut: 表示最少要换出多少代币的数量,0表示让系统看着办(自己算)

最关的是amountIn的代币数量虽然是确定的,但是怎样知道它是这个市场交易对的哪种代币呢?

得去找这个指令对应的内部指令`Inner Instructions`:
这笔交易中有2个对应的内部指令:

`内部指令1`:
```json
{
  "info": {
    "amount": "1000",
    "authority": "vEp6r1ArgChJsKghxCND4TVarEeVo6uSMWsEDFuCKER",
    "destination": "35HseUkS5TMjHeS56NTT6h9PBn1eAJ7STv5fJ2bpit9V",
    "source": "FWrZus8XZfLAn2HtZdrhRxCWoPjMgWRGwZJn2zQ4em73"
  },
  "type": "transfer"
}
```

`内部指令2`:
```json
{
  "info": {
    "amount": "4387433",
    "authority": "5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1",
    "destination": "EoqFJsEyfWvFc4VLWMMiZbuN5aQ4DEymdzhoB3fvixPt",
    "source": "FmJe8E2LssvAU4dSPPoS6iiZTySy7dae114nmDvBTh64"
  },
  "type": "transfer"
}
```

这是从区块获取到的原始数据

要根据内部指令的destination来判断,具体的看`内部指令1`:
- 如果是baseVault,既用户是向baseVault注入代币,说明用户是用base买quote
- 如果是quoteVault,既用户是向quoteVault注入代币,说明用户是用quote买base

那么`内部指令2`一定是从pool1/2转出

所以只看`内部指令1`的destination字段:
- 如果是baseVault那么amountIn对应的代币就是market的baseMint字段
- 如果是quoteVault那么amountIn对应的代币就是market的quoteMint字段

