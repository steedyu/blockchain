package main

import (
	"jerome.com/blockchain/bolt"
	"os"
	"fmt"
)

const dbFile = "blockChain.db"
const blockBucket = "bucket"
const lastHashKey = "Key"
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type BlockChain struct {
	//blocks []*Block 废弃

	//数据库操作句柄
	db   *bolt.DB
	//尾巴，表示最后一个区块的哈希值
	tail []byte
}

func isDBExists() bool {
	// If there is an error, it will be of type *PathError.
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//创建blockchain数据库文件
func InitBlockChain(address string) *BlockChain {
	if isDBExists() {
		fmt.Println("blockchain exist already, no need to create")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("NewBlockChain", err)

	var lastHash []byte

	db.Update(func(tx *bolt.Tx) error {
		coinbase := NewCoinbaseTx(address, genesisInfo)
		//没有bucket,要去创建创世块，将数据天蝎到数据库的bucket中
		genesis := NewGenesisBlock(coinbase)
		bucket, err := tx.CreateBucket([]byte(blockBucket))
		CheckErr("NewBlockChain1", err)
		err = bucket.Put(genesis.Hash, genesis.Serialize())
		CheckErr("NewBlockChain2", err)
		err = bucket.Put([]byte(lastHashKey), genesis.Hash)
		CheckErr("NewBlockChain3", err)
		lastHash = genesis.Hash

		return nil
	})

	return &BlockChain{db, lastHash}
}

func GetBlockChainHandler() *BlockChain {
	if !isDBExists() {
		fmt.Println("Pls create blockchain first")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("GetBlockChainHandler1", err)

	var lastHash []byte

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			//取出最后区块的哈希值
			lastHash = bucket.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

func (bc *BlockChain) AddBlock(txs []*Transaction) {
	var lastHash []byte
	bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	block := NewBlock(txs, lastHash)

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

//返回指定地址的UTXO
func (bc *BlockChain)FindUTXOTransactions(address string) []Transaction {

	//包含目标的utxo的交易集合
	var UTXOTransactions []Transaction = make([]Transaction, 0)
	//存储使用过的utxo的集合 map[交易id] []int64
	//0x011111: 0,1 都是Alice 转账
	spentUTXO := make(map[string][]int64)

	it := bc.NewIterator()
	for {
		//遍历区块
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transactions {
			//遍历input
			//目的：找到已经消耗过的utxo，把他们放到一个集合里
			//需要两个字段来表示使用过的utxo， a 交易ID  b output的索引
			if !tx.IsCoinbase() {
				for _, input := range tx.TXInputs {
					if input.CanUnlockUTXOWith(address) {
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Vout)
					}
				}
			}

			OUTPUTS:
			//遍历output
			//目的：找到所有能支配utxo
			for currentIndex, output := range tx.TXOutputs {
				//检查当前output是否已经被消耗，如果消耗过，那么就进行下一个output检验
				if spentUTXO[string(tx.TXID)] != nil {
					//非空，代表当前交易里面有消耗的utxo
					indexes := spentUTXO[string(tx.TXID)]
					for _, index := range indexes {
						//当前索引河消耗索引比较，有相同，表明这个output肯定被消耗了，直接跳过，进行下一个output判断
						if int64(currentIndex) == index {
							continue OUTPUTS
						}
					}

				}
				//如果当前地址是这个utxo的所有者，就满足条件
				if output.CanBeUnlockWith(address) {
					UTXOTransactions = append(UTXOTransactions, *tx)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return UTXOTransactions
}

//寻找指定地址能够使用的utxo
func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	txs := bc.FindUTXOTransactions(address)

	for _, tx := range txs {
		for _, utxo := range tx.TXOutputs {
			if utxo.CanBeUnlockWith(address) {
				UTXOs = append(UTXOs, utxo)
			}
		}
	}

	return UTXOs
}

func (bc *BlockChain) FindSuitableUTXOs(address string, amount float64) (map[string][]int64, float64) {
	txs := bc.FindUTXOTransactions(address)

	validutxos := make(map[string][]int64)
	var total float64
	FIND:
	//遍历交易
	for _, tx := range txs {
		for index, output := range tx.TXOutputs {
			if output.CanBeUnlockWith(address) {
				//判断当前收集utxo的总金额是否大于所需要花费的金额
				if total < amount {
					validutxos[string(tx.TXID)] = append(validutxos[string(tx.TXID)], int64(index))
					total += output.Value
				} else {
					break FIND
				}
			}
		}
	}

	return validutxos, total
}




















