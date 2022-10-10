package nft

import (
	"context"
	"testing"

	"github.com/coming-chat/go-aptos/aptosclient"
	txnBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/stretchr/testify/require"
)

const (
	testnetRpc = "https://testnet.aptoslabs.com"
)

var (
	restClient, _ = aptosclient.Dial(context.Background(), testnetRpc)
	tokenClient   = NewTokenClient(restClient)
)

func TestGetCollectionData(t *testing.T) {
	address := "0xabf3630d0532fef81dfe610dd4def095070d91e344d475051e1c49da5e6d51c3"
	account, err := txnBuilder.NewAccountAddressFromHex(address)
	require.Nil(t, err)

	data, err := tokenClient.GetCollectionData(*account, "Aptos Zero")
	require.Nil(t, err)
	t.Log(data)
}

func TestGetTokenData(t *testing.T) {
	address := "0xabf3630d0532fef81dfe610dd4def095070d91e344d475051e1c49da5e6d51c3"
	account, err := txnBuilder.NewAccountAddressFromHex(address)
	require.Nil(t, err)

	data, err := tokenClient.GetTokenData(*account, "Aptos Zero", "Aptos Zero: 909587")
	require.Nil(t, err)
	t.Log(data)
}

func TestGetTokenForAccount(t *testing.T) {
	owner, _ := txnBuilder.NewAccountAddressFromHex("0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5")

	address := "0xabf3630d0532fef81dfe610dd4def095070d91e344d475051e1c49da5e6d51c3"
	creator, err := txnBuilder.NewAccountAddressFromHex(address)
	require.Nil(t, err)

	tokenId := TokenId{
		TokenDataId: TokenDataId{
			Creator:    creator.ToShortString(),
			Collection: "Aptos Zero",
			Name:       "Aptos Zero: 909587",
		},
	}
	data, err := tokenClient.GetTokenForAccount(*owner, tokenId)
	require.Nil(t, err)
	t.Log(data)
}

func TestGetAllTokenForAccount(t *testing.T) {
	owner, _ := txnBuilder.NewAccountAddressFromHex("0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5")
	nfts, err := tokenClient.GetAllTokenForAccount(*owner)
	require.Nil(t, err)
	t.Log(nfts)
}
