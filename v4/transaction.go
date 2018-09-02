package main

import (
	"crypto/sha256"
	"encoding/gob"
	"bytes"
	"fmt"
	"os"
)

const reward = 12.5

type Transaction struct {
	//交易ID
	TXID      []byte
	//输入
	TXInputs  []TXInput
	//输出
	TXOutputs []TXOutput
}

func (tx *Transaction) IsCoinbase() bool {
	if len(tx.TXInputs) == 1 {
		if len(tx.TXInputs[0].TXID) == 0 && tx.TXInputs[0].Vout == -1 {
			return true
		}
	}

	return false
}

//创建coinbase交易，只有收款人，没有付款人,是矿工的奖励交易
func NewCoinbaseTx(address string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("reward to %s %d btc", address, reward)
	}

	input := TXInput{TXID:[]byte{}, Vout:-1, ScriptSig:data}
	output := TXOutput{Value:reward, ScriptPubKey: address}

	tx := Transaction{TXID:[]byte{}, TXInputs: []TXInput{input}, TXOutputs:[]TXOutput{output} }
	tx.SetTXID()
	return &tx
}

//创建一个普通交易，send的辅助函数
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	//map[string][]int64 key : 交易id， value: 引用output的索引数据
	var validUTXOs map[string][]int64 = make(map[string][]int64)  /*所需要的合理的UTXO */
	var total float64 /*返回UTXO的金额总和 */
	validUTXOs, total = bc.FindSuitableUTXOs(from, amount)

	/*
	validUTXOs[0x11111111] = []int64{1}
	validUTXOs[0x22222222] = []int64{0}
	....
	validUTXOs[0xnnnnnnn] = []int64{0,4,8}
	*/
	if total < amount {
		fmt.Println("Not Enough Money!")
		os.Exit(1)
	}

	var inputs []TXInput
	var outputs []TXOutput

	//1 创建inputs
	//进行output到input转换
	//遍历有效utxo的合集
	for txId, outputIndexes := range validUTXOs {
		//遍历所有引用的utxo索引，每一个索引需要创建一个input
		for _, index := range outputIndexes {
			input := TXInput{
				TXID : []byte(txId),
				Vout : index,
				ScriptSig : from,
			}
			inputs = append(inputs, input)
		}
	}

	//2 创建outputs
	//给对方支付
	output := TXOutput{Value:amount, ScriptPubKey:to}
	outputs = append(outputs, output)

	//找零钱
	if total > amount {
		output := TXOutput{Value:total - amount, ScriptPubKey:from}
		outputs = append(outputs, output)
	}

	tx := Transaction{TXID:nil, TXInputs: inputs, TXOutputs: outputs }
	tx.SetTXID()
	return &tx
}

type TXInput struct {
	//所引用输出的交易ID
	TXID      []byte
	//所引用output的索引值
	Vout      int64
	//解锁脚本，指明可以使用某个output的条件
	ScriptSig string
}

//判断下你提供的信息，是否能解锁使用这个UTXO
//检查当前用户能否借口引用的UTXO
func (input *TXInput) CanUnlockUTXOWith(unlockData string) bool {
	return unlockData == input.ScriptSig
}

type TXOutput struct {
	//支付给收款方的金额
	Value        float64
	//锁定脚本，制定收款方的地址
	ScriptPubKey string
}

//所查到的UTXO 是否能被我的私钥所解开
//检查当前用户是否是这个utxo的所有者
func (output *TXOutput) CanBeUnlockWith(unlockData string) bool {
	return unlockData == output.ScriptPubKey
}

func (tx *Transaction) SetTXID() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)

	CheckErr("Serialize", err)
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}