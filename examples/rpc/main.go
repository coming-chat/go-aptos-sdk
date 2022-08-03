package main

import (
	"context"
	"encoding/json"
	"fmt"

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
}
