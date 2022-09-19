package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
)

const rpcUrl = "https://fullnode.devnet.aptoslabs.com"

func main() {
	ctx := context.Background()
	client, err := aptosclient.Dial(ctx, rpcUrl)
	if err != nil {
		printError(err)
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		panic(err)
	}
	content, err := json.Marshal(ledgerInfo)
	if err != nil {
		printError(err)
	}
	fmt.Println(string(content))

	transactions, err := client.GetTransactions(1, 10)
	if err != nil {
		printError(err)
	}
	printLine("get tx list")
	for _, tx := range transactions {
		fmt.Printf("version: %d, type: %s, hash: %s, round: %d, time: %d\n", tx.Version, tx.Type, tx.Hash, tx.Round, tx.Timestamp)
	}

	events, err := client.GetEventsByCreationNumber("0x8ab4fe552d76c49cdf7a4707338520ddf5de705044665fbcc4b6ea60a6b026d4", "2", 0, 10)
	if err != nil {
		printError(err)
	}
	for _, e := range events {
		fmt.Println(e.SequenceNumber)
	}
}

func printLine(content string) {
	fmt.Printf("================= %s =================\n", content)
}

func printError(err error) {
	var restError *aptostypes.RestError
	if b := errors.As(err, &restError); b {
		fmt.Printf("code: %d, message: %s, aptos_ledger_version: %d\n", restError.Code, restError.Message, restError.AptosLedgerVersion)
	} else {
		fmt.Printf("err: %s\n", err.Error())
	}
}
