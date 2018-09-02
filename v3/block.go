package main

import (
	"time"
	"bytes"
	"encoding/gob"
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

func (block *Block) Serialize() []byte {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckErr("Serialize", err)

	return buffer.Bytes()
}

func Deserialize(data []byte) *Block {
	if len(data) == 0 {
		return nil
	}

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	CheckErr("Deserialize", err)
	return &block
}

func NewBlock(data string, prevBlockhash []byte) *Block {

	var block Block
	block = Block{
		Version:1,
		PrevBlockHash:prevBlockhash,
		MerKelRoot:[]byte{},
		TimeStamp:time.Now().Unix(),
		Bits:targetBits,
		Nonce:0,
		Data: []byte(data),
	}

	//block.SetHash()
	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return &block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block!", []byte{})
}







