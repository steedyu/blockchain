package main

/*
v4版本思路

1 将创建区块链的操作放到命令
   NewBlockChain

2 定义交易结构
 a 交易ID
 b 交易输入:TXInput
 c 交易输出:TXOutput

3 根据交易结构，改写代码
  a 创建区块链的时候生成奖励
  b 通过指定地址检索到他相关的UTX0
  c 实现UTXO的转移(创建交易函数:NewTransaction (from,to string, amount float64))

4 实现命令
send --from From --to TO --amount AMOUNT 	"send coin from From to TO"
	getbalance --address ADRESS	"get balance of the address"
 */

func main() {

	//cmd CreateChain, getBalance

	cli := CLI{}
	cli.Run()
}
