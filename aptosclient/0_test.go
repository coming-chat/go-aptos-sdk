package aptosclient

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	MainnetRestUrl = "https://fullnode.mainnet.aptoslabs.com"
	TestnetRestUrl = "https://testnet.aptoslabs.com"
	DevnetRestUrl  = "https://fullnode.devnet.aptoslabs.com"
)

func Client(t *testing.T, url string) RestClient {
	cli, err := Dial(context.Background(), url)
	require.Nil(t, err)
	return *cli
}

func txnSubmitableForTest(t *testing.T) bool {
	out, _ := exec.Command("whoami").Output()
	user := strings.TrimSpace(string(out))
	switch user {
	case "ggggg":
		return true
	default:
		t.Log("Non-specified machines, stop sending transactions after signing: ", user)
		return false
	}
}
