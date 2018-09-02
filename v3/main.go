package main

/*
V3版本思路
1 bolt数据库介绍
	key-〉value进行读取，存储
	轻量级
	开源的

2 NewBlockChain函数重写
	由数组编程数据库操作
	创建数据库文件

3 AddBlock函数重写
	对数据的读取和写入

4 对数据库的遍历
	迭代器的编写，Iterator

5 命令行介绍及编写
	a 添加区块命令
	b 打印区块链命令
 */

func main() {
	bc := NewBlockChain()
	//bc.AddBlock("A send B 1BTC")
	//bc.AddBlock("B send C 1BTC")

	//for _, block := range bc.blocks {
	//	fmt.Printf("Version: %dd\n", block.Version)
	//	fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
	//	fmt.Printf("Hash: %x\n", block.Hash)
	//	fmt.Printf("TimeStamp: %d\n", block.TimeStamp)
	//	fmt.Printf("Bits: %d\n", block.Bits)
	//	fmt.Printf("Nonce: %d\n", block.Nonce)
	//	fmt.Printf("Data: %s\n", block.Data)
	//	fmt.Printf("IsValid: %v\n", NewProofOfWork(block).IsValid())
	//}

	cli := CLI{bc:bc}
	cli.Run()
}
