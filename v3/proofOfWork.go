package main

import (
	"math/big"
	"crypto/sha256"
	"bytes"
	"math"
	"fmt"
)

type ProofOfWork struct {
	block  *Block

	//目标值
	target *big.Int
}

const targetBits = 24

func NewProofOfWork(block *Block) *ProofOfWork {

	//000000000000000......01
	target := big.NewInt(1)
	//0x00000100000000000000
	target.Lsh(target, uint(256 - targetBits))
	pow := ProofOfWork{block:block, target:target}
	return &pow
}

func (pow *ProofOfWork) PrepareData(nonce int64) []byte {

	block := pow.block
	tmp := [][]byte{
		IntToBytes(block.Version),
		block.PrevBlockHash,
		block.MerKelRoot,
		IntToBytes(block.TimeStamp),
		IntToBytes(targetBits),
		IntToBytes(nonce),
		block.Data,
	}
	data := bytes.Join(tmp, []byte{})
	return data
}

func (pow *ProofOfWork) Run() (int64, []byte) {

	//1 拼装数据
	//2 哈希值转成big.Int
	/*
	for nonce {
		hash := sha256(block数据 + nonce)
		if 转换(hash) < pow.target {
			找到了
		}else {
			nonce ++
		}
	}
	return nonce hash[:]
	*/

	var hash [32]byte
	var nonce int64 = 0
	var hashInt big.Int

	fmt.Println("Begin Mining...")
	fmt.Printf("target hash :   %x\n", pow.target.Bytes())
	for nonce < math.MaxInt64 {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("found nonce, nonce :%d, hash:%x\n", nonce, hash)
			break
		} else {
			nonce ++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}