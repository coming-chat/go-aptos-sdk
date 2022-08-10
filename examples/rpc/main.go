package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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

	account := "0x68fd7b5e581d2a95c5dbba09b9c19879c7a934f4f4cfda1d5008cc793660c8ee"
	accountTransactions, err := client.GetAccountTransactions(account, 1, 10)
	if err != nil {
		printError(err)
	}
	printLine("get account tx list")
	for _, tx := range accountTransactions {
		fmt.Printf("type: %s, hash: %s, version: %d\n", tx.Type, tx.Hash, tx.Version)
	}

	txHash := "0xe1163f7b33df37995c9724b179b4d4a0fdff9eb2ef0c38a8d6e2982ce1c1de22"
	tx, err := client.GetTransaction(txHash)
	if err != nil {
		printError(err)
	}
	printLine("get tx by hash")
	fmt.Printf("type: %s, hash: %s, version: %d\n", tx.Type, tx.Hash, tx.Version)

	if tx.Type == aptostypes.TypeUserTransaction {
		userTx := tx.AsUserTransaction()
		fmt.Printf("type: %s, hash: %s, version: %d\n", userTx.Type, userTx.Hash, userTx.Version)
	}

	tx, err = client.GetTransaction(strconv.FormatUint(tx.Version, 10))
	if err != nil {
		printError(err)
	}
	printLine("get tx by version")
	fmt.Printf("type: %s, hash: %s, version: %d\n", tx.Type, tx.Hash, tx.Version)

	tx, err = client.GetTransaction(strconv.FormatUint(6618578, 10))
	if err != nil {
		printError(err)
	}
	if tx.Type == aptostypes.TypeBlockMetadataTransaction {
		blockTx := tx.AsBlockMetadataTransaction()
		fmt.Printf("type: %s, hash: %s, version: %d, id: %s\n", blockTx.Type, blockTx.Hash, blockTx.Version, blockTx.ID)
	}

	accountCore, err := client.GetAccount(account)
	if err != nil {
		printError(err)
	}
	fmt.Printf("seqNum: %d, key: %s\n", accountCore.SequenceNumber, accountCore.AuthenticationKey)

	accountResources, err := client.GetAccountResources(account)
	if err != nil {
		printError(err)
	}
	printLine("account resource")
	for _, resource := range accountResources {
		fmt.Printf("resourceType: %s, data: %v\n", resource.Type, resource.Data)
	}

	resourceType := "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
	resource, err := client.GetAccountResource(account, resourceType, 0)
	if err != nil {
		printError(err)
	}
	fmt.Printf("resourceType: %s\n, data: %v\n", resource.Type, resource.Data)

	accountWithModule := "0xe72ca8ac40bb39bf4cf6d5abf8655015e3fce0f0a5e4218935355625eb0224a3"
	accountModules, err := client.GetAccountModules(accountWithModule, 0)
	if err != nil {
		printError(err)
	}
	printLine("account modules")
	for _, module := range accountModules {
		fmt.Printf("abi:%s, name:%s\n", module.Abi.Address, module.Abi.Name)
	}

	accountModule, err := client.GetAccountModule(accountWithModule, "message", 0)
	printLine("account module")
	if err != nil {
		printError(err)
	} else {
		fmt.Printf("abi:%s, name:%s\n", accountModule.Abi.Address, accountModule.Abi.Name)
	}

	eventKey := "0x0100000000000000874342f90ed0c0ccdf7baa13309820133ef94f143bb4a68069ceae8a8658541a"
	events, err := client.GetEventsByKey(eventKey)
	if err != nil {
		printError(err)
	}
	printLine("events by key")
	for _, event := range events {
		fmt.Printf("key: %s, seqNum: %d, type: %s\n", event.Key, event.SequenceNumber, event.Type)
	}

	addressWithEvent := "0x647040d2018e65ae91a2353125a06a7b58917c523bef4e775237a814e464918c"
	eventHandleStruct := "0x647040d2018e65ae91a2353125a06a7b58917c523bef4e775237a814e464918c::message::MessageHolder"
	fieldName := "message_change_events"
	events, err = client.GetEventsByEventHandle(addressWithEvent, eventHandleStruct, fieldName, 0, 0)
	if err != nil {
		printError(err)
	}
	printLine("events by address/handle/field")
	for _, event := range events {
		fmt.Printf("key: %s, seqNum: %d, type: %s\n", event.Key, event.SequenceNumber, event.Type)
	}

	printLine("error test")
	_, err = client.GetEventsByEventHandle("0x123", eventHandleStruct, fieldName, 0, 0)
	if err != nil {
		printError(err)
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
