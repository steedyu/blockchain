package main

import (
	"os"
	"fmt"
	"flag"
)

const usage = `
	createChain --address ADDRESS "create a blockchain"
	send --from From --to TO --amount AMOUNT 	"send coin from From to TO"
	getBalance --address ADRESS	"get balance of the address"
	printChain		"print all block"
`

//const AddBlockCmdString = "addBlock"
const PrintChainCmdString = "printChain"
const CreateChainCmdString = "createChain"
const getBalanceCmdString = "getBalance"
const sendCmdString = "send"

type CLI struct {

}

func (cli *CLI) printUsage() {
	fmt.Println(usage)
	os.Exit(1)
}

func (cli *CLI) parameterCheck() {
	if len(os.Args) < 2 {
		fmt.Println("invalid input!")
		cli.printUsage()
	}
}

func (cli *CLI) Run() {
	cli.parameterCheck()

	//addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	createChainCmd := flag.NewFlagSet(CreateChainCmdString, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(getBalanceCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(sendCmdString, flag.ExitOnError)

	//addblockCmdPara := addBlockCmd.String("data", "", "block transaction info!")
	//创建区块链相关参数
	createChainCmdPara := createChainCmd.String("address", "", "address info!")
	//余额相关参数
	getBalanceCmdPara := getBalanceCmd.String("address", "", "address info!")
	//相关参数
	fromPara := sendCmd.String("from", "", "sender address info!")
	toPara := sendCmd.String("to", "", "to address info!")
	amountPara := sendCmd.Float64("amount", 0, "amount info!")

	switch os.Args[1] {
	//case AddBlockCmdString:
	//	err := addBlockCmd.Parse(os.Args[2:])
	//	CheckErr("Run()", err)
	//	if addBlockCmd.Parsed() {
	//		if *addblockCmdPara == "" {
	//			cli.printUsage()
	//		}
	//
	//		cli.AddBlock(*addblockCmdPara)
	//	}
	case CreateChainCmdString:
		err := createChainCmd.Parse(os.Args[2:])
		CheckErr("Run3()", err)
		if createChainCmd.Parsed() {
			if *createChainCmdPara == "" {
				cli.printUsage()
			}

			cli.CreateChain(*createChainCmdPara)
		}
	case getBalanceCmdString:
		err := getBalanceCmd.Parse(os.Args[2:])
		CheckErr("Run4()", err)
		if getBalanceCmd.Parsed() {
			if *getBalanceCmdPara == "" {
				cli.printUsage()
			}

			cli.GetBalance(*getBalanceCmdPara)
		}
	case sendCmdString:
		err := sendCmd.Parse(os.Args[2:])
		CheckErr("Run5()", err)
		if sendCmd.Parsed() {
			if *fromPara == "" || *toPara == "" || *amountPara <= 0 {
				cli.printUsage()
			}

			cli.Send(*fromPara, *toPara, *amountPara)
		}

	case PrintChainCmdString:
		err := printChainCmd.Parse(os.Args[2:])
		CheckErr("Run2()", err)
		if printChainCmd.Parsed() {
			cli.PrintChain()
		}
	default:
		cli.printUsage()
	}
}
