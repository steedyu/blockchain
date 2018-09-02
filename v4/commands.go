package main

import "fmt"

//命令的实现放到另外一个文件便于模块化管理

func (cli *CLI) AddBlock(data string) {
	//bc := GetBlockChainHandler()
	//bc.AddBlock(data)
}

func (cli *CLI) PrintChain() {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	it := bc.NewIterator()
	for {
		block := it.Next()

		fmt.Printf("Version: %dd\n", block.Version)
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("TimeStamp: %d\n", block.TimeStamp)
		fmt.Printf("Bits: %d\n", block.Bits)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		//fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("IsValid: %v\n", NewProofOfWork(block).IsValid())

		if len(block.PrevBlockHash) == 0 {
			fmt.Println("print over!")
			break
		}
	}
}

func (cli *CLI) CreateChain(address string) {
	bc := InitBlockChain(address)
	defer bc.db.Close()
	fmt.Println("Init blockchain successfully!")
}

func (cli *CLI) GetBalance(address string) {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	utxoes := bc.FindUTXO(address)

	//总金额
	var total float64
	//遍历所有的utxo，获取金额总数
	for _, utxo := range utxoes {
		total += utxo.Value
	}

	fmt.Printf("The balance of %s is %f \n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64) {
	bc := GetBlockChainHandler()
	defer bc.db.Close()

	tx := NewTransaction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{tx})
	fmt.Println("send successfully!")
}

