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
		fmt.Printf("resourceType: %s\n, data: %v", resource.Type, resource.Data)
	}

	resourceType := "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
	resource, err := client.GetAccountResource(account, resourceType, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("resourceType: %s\n, data: %v", resource.Type, resource.Data)

}

func printLine(content string) {
	fmt.Printf("================= %s =================\n", content)
}
