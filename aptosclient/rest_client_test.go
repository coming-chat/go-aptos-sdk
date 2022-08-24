package aptosclient

import (
	"context"
	"testing"
)

// const RestUrl = "https://aptosdev.coming.chat/v1"
const RestUrl = "https://fullnode.devnet.aptoslabs.com/"

func TestLedgerInfo(t *testing.T) {
	client, err := Dial(context.Background(), RestUrl)
	if err != nil {
		t.Fatal(err)
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ledgerInfo)
}

func TestTransactionDetail(t *testing.T) {
	client, err := Dial(context.Background(), RestUrl)
	if err != nil {
		t.Fatal(err)
	}

	version := "16603388"
	tx1, err := client.GetTransactionByVersion(version)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx.hash = %v", tx1.Hash)

	hash := "0x7e2eabcf1dfd252599f8ae1369a945c9fda947bdffb6e35d37873beed2463ddb"
	tx2, err := client.GetTransactionByHash(hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx.version = %v", tx2.Version)
}
