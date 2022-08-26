package aptosclient

import (
	"context"
	"strconv"
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

	version := "500000"
	tx1, err := client.GetTransactionByVersion(version)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx.hash = %v", tx1.Hash)

	hash := tx1.Hash
	tx2, err := client.GetTransactionByHash(hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx.version = %v", tx2.Version)

	if strconv.FormatUint(tx2.Version, 10) != version {
		t.Fatal("Transaction's version and hash not match.")
	}
}
