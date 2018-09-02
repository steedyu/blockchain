package main

import "fmt"

/*

v2版本
1 定义一个工作量证明的结构ProofOfWork
	block
	目标值

2 创建一个POW的方法
	NewProofOfWork(参数)

3 提供一个计算哈希值方法
	Run()
4 提供一个校验函数
	isValid()



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

	}
}
