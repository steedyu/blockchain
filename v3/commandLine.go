package main

import (
	"os"
	"fmt"
	"flag"
)

const usage = `
	addBlock --data DATA	"add a block to blockchain"
	printChain		"print all block"
`

const AddBlockCmdString = "addBlock"
const PrintChainCmdString = "printChain"

type CLI struct {
	bc *BlockChain
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

	addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)

	addblockCmdPara := addBlockCmd.String("data", "", "block transaction info!")

	switch os.Args[1] {
	case AddBlockCmdString:
		err := addBlockCmd.Parse(os.Args[2:])
		CheckErr("Run()", err)
		if addBlockCmd.Parsed() {
			if *addblockCmdPara == "" {
				cli.printUsage()
			}

			cli.AddBlock(*addblockCmdPara)
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
