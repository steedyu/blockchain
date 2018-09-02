package main

import (
	"time"
	"crypto/sha256"
	"bytes"
)

type Block struct {
	//版本
	Version       int64
	//前区块哈希值
	PrevBlockHash []byte
	//梅克尔根
	MerKelRoot    []byte
	//为了简便，这里就不计算当前区块hash值，而是保存在里面
	Hash          []byte
	//时间戳
	TimeStamp     int64
	//难度值
	Bits          int64
	//随机值
	Nonce         int64

	//交易信息
	Data          []byte
}

func (block *Block) SetHash() {

	tmp := [][]byte{
		IntToBytes(block.Version),
		block.PrevBlockHash,
		block.MerKelRoot,
		IntToBytes(block.TimeStamp),
		IntToBytes(block.Bits),
		IntToBytes(block.Nonce),
		block.Data,
	}
	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

func NewBlock(data string, prevBlockhash []byte) *Block {

	var block Block
	block = Block{
		Version:1,
		PrevBlockHash:prevBlockhash,
		MerKelRoot:[]byte{},
		TimeStamp:time.Now().Unix(),
		Bits:1, //是根据现有区块的创建时间 计算出来的，简便写1
		Nonce:1, //随机值应该是算出来，简便写1
		Data: []byte(data),
	}

	block.SetHash()
	return &block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block!", []byte{})
}





