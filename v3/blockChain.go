package main

import (
	"jerome.com/blockchain/bolt"
	"os"
)

const dbFile = "blockChain.db"
const blockBucket = "bucket"
const lastHashKey = "Key"

type BlockChain struct {
	//blocks []*Block 废弃

	//数据库操作句柄
	db   *bolt.DB
	//尾巴，表示最后一个区块的哈希值
	tail []byte
}

func NewBlockChain() *BlockChain {

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("NewBlockChain", err)

	var lastHash []byte

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			//取出最后区块的哈希值
			lastHash = bucket.Get([]byte(lastHashKey))
		} else {
			//没有bucket,要去创建创世块，将数据天蝎到数据库的bucket中
			genesis := NewGenesisBlock()
			bucket, err := tx.CreateBucket([]byte(blockBucket))
			CheckErr("NewBlockChain2", err)
			err = bucket.Put(genesis.Hash, genesis.Serialize())
			CheckErr("NewBlockChain3", err)
			err = bucket.Put([]byte(lastHashKey), genesis.Hash)
			CheckErr("NewBlockChain4", err)
			lastHash = genesis.Hash
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte
	bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	block := NewBlock(data, lastHash)

	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		err := bucket.Put(block.Hash, block.Serialize())
		CheckErr("AddBlock1", err)
		err = bucket.Put([]byte(lastHashKey), block.Hash)
		CheckErr("AddBlock2", err)
		bc.tail = block.Hash
		return nil
	})
	CheckErr("AddBlock3", err)
}

//迭代器 就是一个对象 它里面包含一个游标 一直向前（后）移动，完成整个容器的遍历
type BlockChainIterator struct {
	currHash []byte
	db       *bolt.DB
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		currHash:bc.tail,
		db:bc.db,
	}
}

func (it *BlockChainIterator) Next() (block *Block) {
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			return nil
		}

		data := bucket.Get(it.currHash)
		block = Deserialize(data)
		it.currHash = block.PrevBlockHash
		return nil
	})
	CheckErr("Next()", err)
	return
}