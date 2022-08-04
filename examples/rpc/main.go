package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/coming-chat/go-aptos/aptosclient"
)

const rpcUrl = "https://fullnode.devnet.aptoslabs.com"

func main() {
	ctx := context.Background()
	client, err := aptosclient.Dial(ctx, rpcUrl)
	if err != nil {
		panic(err)
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		panic(err)
	}
	content, err := json.Marshal(ledgerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))

	transactions, err := client.GetTransactions(ledgerInfo.LedgerVersion-10, 10)
	if err != nil {
		panic(err)
	}
	printLine("get tx list")
	for _, tx := range transactions {
		fmt.Printf("type: %s, hash: %s\n", tx.Type, tx.Hash)
	}

	account := "0xa1f475d2190bb689fa68804bb0be954c640d582290fbb49aa05c4d438c989603"
	accountTransactions, err := client.GetAccountTransactions(account, 1, 10)
	if err != nil {
		panic(err)
	}
	printLine("get account tx list")
	for _, tx := range accountTransactions {
		fmt.Printf("type: %s, hash: %s, version: %d\n", tx.Type, tx.Hash, tx.Version)
	}

	txHash := "0x95f1df00314e740a71acb66aeb7dff0182f7510bf900c356db91318b8952ed1d"
	tx, err := client.GetTransaction(txHash)
	if err != nil {
		panic(err)
	}
	printLine("get tx by hash")
	fmt.Printf("type: %s, hash: %s, version: %d\n", tx.Type, tx.Hash, tx.Version)

	tx, err = client.GetTransaction(strconv.FormatUint(6047729, 10))
	if err != nil {
		panic(err)
	}
	printLine("get tx by version")
	fmt.Printf("type: %s, hash: %s, version: %d\n", tx.Type, tx.Hash, tx.Version)

	accountCore, err := client.GetAccount(account)
	if err != nil {
		panic(err)
	}
	fmt.Printf("seqNum: %d, key: %s\n", accountCore.SequenceNumber, accountCore.AuthenticationKey)

	accountResources, err := client.GetAccountResources(account)
	if err != nil {
		panic(err)
	}
	printLine("account resource")
	for _, resource := range accountResources {
		fmt.Printf("resourceType: %s, data: %v\n", resource.Type, resource.Data)
	}

	resourceType := "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
	resource, err := client.GetAccountResource(account, resourceType, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("resourceType: %s\n, data: %v\n", resource.Type, resource.Data)

	accountWithModule := "0xe72ca8ac40bb39bf4cf6d5abf8655015e3fce0f0a5e4218935355625eb0224a3"
	accountModules, err := client.GetAccountModules(accountWithModule, 0)
	if err != nil {
		panic(err)
	}
	printLine("account modules")
	for _, module := range accountModules {
		fmt.Printf("abi:%s, name:%s\n", module.Abi.Address, module.Abi.Name)
	}

	accountModule, err := client.GetAccountModule(accountWithModule, "message", 0)
	if err != nil {
		panic(err)
	}
	printLine("account module")
	fmt.Printf("abi:%s, name:%s\n", accountModule.Abi.Address, accountModule.Abi.Name)

	eventKey := "0x0100000000000000874342f90ed0c0ccdf7baa13309820133ef94f143bb4a68069ceae8a8658541a"
	events, err := client.GetEventsByKey(eventKey)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	printLine("events by address/handle/field")
	for _, event := range events {
		fmt.Printf("key: %s, seqNum: %d, type: %s\n", event.Key, event.SequenceNumber, event.Type)
	}
}

func printLine(content string) {
	fmt.Printf("================= %s =================\n", content)
}
