package aptosclient

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLedgerInfo(t *testing.T) {
	client := Client(t, DevnetRestUrl)
	ledgerInfo, err := client.LedgerInfo()
	require.Nil(t, err)
	t.Log(ledgerInfo)
}

func TestTransactionDetail(t *testing.T) {
	client := Client(t, DevnetRestUrl)

	version := "500000"
	tx1, err := client.GetTransactionByVersion(version)
	require.Nil(t, err)
	t.Logf("tx.hash = %v", tx1.Hash)

	hash := tx1.Hash
	tx2, err := client.GetTransactionByHash(hash)
	require.Nil(t, err)
	t.Logf("tx.version = %v", tx2.Version)

	tx2Version := strconv.FormatUint(tx2.Version, 10)
	require.Equal(t, tx2Version, version, "Transaction's version and hash not match.")
	require.Equal(t, tx1.Hash, tx2.Hash)
}
