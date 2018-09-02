package main

import "fmt"

/*

v1版本思路

区块相关
1 定义一个区块的结构Block
	a.区块头 6个字段
	b.区块体：字符串表示data

2 提供一个创建区块的方法
	NewBlock(参数)


区块链相关
1  定义一个区块链结构BlockChain
	Block数组

2 提供一个创建BlockChain的方法
	NewBlockChain()

3 提供一个添加区块的方法
	AddBlock(参数)


 */

func main() {
	bc := NewBlockChain()
	bc.AddBlock("A send B 1BTC")
	bc.AddBlock("B send C 1BTC")

	for _, block := range bc.blocks {
		fmt.Printf("Version: %dd\n", block.Version)
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("TimeStamp: %d\n", block.TimeStamp)
		fmt.Printf("Bits: %d\n", block.Bits)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("IsValid: %v\n", NewProofOfWork(block).IsValid())
	}
}
