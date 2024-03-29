# neotest

一款测试NEO交易的自动化测试工具



## 构造交易

* UTXO转账
* 合约调用
* 签名
* 支付手续费
* 自定义执行脚本

## NTF语言手册

### 命令

#### echo

都用于调试、打印信息。语法：

```bash
echo <object1> <object2> ...
```

#### tx

发起一次交易， 接收一个交易描述

```bash
tx "transfer-neo"
```

#### tx-v

指定交易版本

```bash
tx-v 0|1
```

#### tx-type

指定交易类型，目前只支持`ContractTransaction`和`InvocationTransaction`

```bash
tx-type "contract"|"invocation"
```

#### tx-fee

指定网络费用

```bash
te-fee <float>
```

#### tx-attr

指定交易的Attribute，接收两个字符串参数，其中第一个参数为usage，值有ContractHash, ECDH02, ECDH03, Script, Vote, DescriptionUrl, Description, Hash1-Hash15, Remark, Remark1-Remark15；第二个参数为data。

该命令可以声明多个。

```bash
tx-attr "Script" "5e40b22e86dc6ff4a7b0416450971469fe71040d"
tx-attr "Hash1" "abcd"
```

#### tx-initiator

交易发起者，接收一个私钥，用于转账UTXO模型代币、支付手续费和交易签名. 构建交易时会使用该私钥对交易进行签名并放入交易的`witness`中

```bash
tx-initiator <PrivateKey>
```

#### tx-vout

指定UTXO模型代币转账, 接收三个参数：

1. 资产哈希，字符串类型
2. 钱包地址，字符串
3. 转账数量，数量值类型

该命令可以多次声明

```bash
tx-vout <asset_hash> <address> <value>
```

#### tx-invoke

使用给定的参数以散列值调用智能合约，接收一个json数组，方便构造自定义脚本。数组元素为构建执行脚本参数。
包含`type`和`value`两个字段；`type`支持以下类型：

* Boolean
* Integer
* Hash160
* Hash256
* String
* Array
* AppCall
* Address
* OpCode

```bash
tx-invoke '[
[
  {
    "type": "AppCall",
    "value": "dc675afc61a7c0f7b3d2682bf6e1d8ed865a0e5f"
  },
  {
    "type": "String",
    "value": "name"
  },
  {
    "type":"Boolean",
    "value": false
  }
]'
```

#### tx-invokefunc

使用给定的操作和参数，以散列值调用智能合约。 接收json数组。

scripthash：智能合约脚本散列值。

operation：操作名称（字符串）。

params：传递给智能合约操作的参数。
  `type`字段支持以下类型：
* Boolean
* Integer
* Hash160
* Hash256
* String
* Array
* AppCall
* Address
* OpCode

```bash
tx-invokefunc '[
    "af7c7328eee5a275a3bcaee2bf0cf662b5e739be",
    "balanceOf",
    [
      {
        "type": "Hash160",
        "value": "91b83e96f2a7c4fdf0c1688441ec61986c7cae26"
      }
    ]
  ]'
```

#### tx-invokescript

指定调用的脚本，接收脚本的十六进制的字符串。

```bash
tx-invokesscript <Script>
```

#### tx-witness

指定见证人，接收两个字符串参数。其中

1. 普通地址的私钥或智能合约的哈希脚本。

2. 调用脚本，当witness为只能合约的哈希脚本时，该值为验证智能合约的参数；当witness为普通地址的私钥时，该值为空即可。

该命令可以多次声明使用

```bash
tx-witness <witness> <invocation>
```

#### tx-send

广播交易。接收一个节点地址

```bash
tx-send "http://localhost:20332"
```



#### req

发起一个HTTP请求

```bash
req <http-method> <url>
```

#### body

 用于指定请求包或返回包的`json`格式正文内容。语法：

```bash
body <json-data>
```

#### ret

ret 用来获得返回包。语法：

```bash
ret [<status-code>]
```

#### let

let 用于变量赋值，和主流命令式编程语言的 `=` 最为接近。例如：

```bash
let $(var-name) <expression>
```

#### equal

equal 要求两个 object 的内容精确相等：

```bash
equal <object1> <object2>
```

### 子命令

#### env

取得环境变量的值

#### now

返回当前时间

