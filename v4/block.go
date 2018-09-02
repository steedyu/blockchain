package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"crypto/sha256"
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
	//Data          []byte
	Transactions  []*Transaction
}

func (block *Block) Serialize() []byte {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckErr("Serialize", err)

	return buffer.Bytes()
}

//粗略模拟梅克尔树，将交易的哈希值进行拼接，生成root hash
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	//遍历交易
	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.TXID)
	}

	//对二维切片进行拼接，生成一维切片
	data := bytes.Join(txHashes, []byte{})
	hash := sha256.Sum256(data)
	return hash[:]
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

func NewBlock(txs []*Transaction, prevBlockhash []byte) *Block {

	var block Block
	block = Block{
		Version:1,
		PrevBlockHash:prevBlockhash,
		MerKelRoot:[]byte{},
		TimeStamp:time.Now().Unix(),
		Bits:targetBits,
		Nonce:0,
		Transactions: txs,
	}

	//block.SetHash()
	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return &block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}








