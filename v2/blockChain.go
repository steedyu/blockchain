package main

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	block := NewGenesisBlock()
	return &BlockChain{blocks: []*Block{block}}
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlockHash := bc.blocks[len(bc.blocks) - 1].Hash
	block := NewBlock(data, prevBlockHash)
	bc.blocks = append(bc.blocks, block)
}