package aptosclient

import (
	"context"
	"testing"
)

// const RestUrl = "https://aptosdev.coming.chat/v1"
const RestUrl = "https://fullnode.devnet.aptoslabs.com/v1"

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
