package nft

import (
	"context"
	"testing"

	"github.com/coming-chat/go-aptos/aptosclient"
	txnBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/stretchr/testify/require"
)

const (
	MainnetRestUrl = "https://fullnode.mainnet.aptoslabs.com"
	TestnetRestUrl = "https://testnet.aptoslabs.com"
	DevnetRestUrl  = "https://fullnode.devnet.aptoslabs.com"
)

const (
	nftCreator = "0x305a97874974fdb9a7ba59dc7cab7714c8e8e00004ac887b6e348496e1981838"
	nftOwner   = "0xa6de5dd7668e3c7b1d8b9d3679b12ed937423fdec5198524665fef3f4799498f"

	nftCollectionName = "Aptos Names V1"
	nftTokenNameOwned = "linmor.apt"
)

var (
	restClient, _ = aptosclient.Dial(context.Background(), MainnetRestUrl)
	tokenClient   = NewTokenClient(restClient)
)

func TestGetCollectionData(t *testing.T) {
	account, err := txnBuilder.NewAccountAddressFromHex(nftCreator)
	require.Nil(t, err)

	data, err := tokenClient.GetCollectionData(*account, nftCollectionName)
	require.Nil(t, err)
	t.Log(data)
}

func TestGetTokenData(t *testing.T) {
	account, err := txnBuilder.NewAccountAddressFromHex(nftCreator)
	require.Nil(t, err)

	data, err := tokenClient.GetTokenData(*account, nftCollectionName, nftTokenNameOwned)
	require.Nil(t, err)
	t.Log(data)
}

func TestGetTokenForAccount(t *testing.T) {
	owner, err := txnBuilder.NewAccountAddressFromHex(nftOwner)
	require.Nil(t, err)
	creator, err := txnBuilder.NewAccountAddressFromHex(nftCreator)
	require.Nil(t, err)

	tokenId := TokenId{
		TokenDataId: TokenDataId{
			Creator:    creator.ToShortString(),
			Collection: nftCollectionName,
			Name:       nftTokenNameOwned,
		},
	}
	data, err := tokenClient.GetTokenForAccount(*owner, tokenId)
	require.Nil(t, err)
	t.Log(data)
}

func TestGetAllTokenForAccount(t *testing.T) {
	owner, err := txnBuilder.NewAccountAddressFromHex(nftOwner)
	require.Nil(t, err)
	nfts, err := tokenClient.GetAllTokenForAccount(*owner)
	require.Nil(t, err)
	t.Log(nfts)
}
