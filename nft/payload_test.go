package nft

import (
	"testing"

	txnBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/stretchr/testify/require"
)

var nftBuilder, _ = NewNFTPayloadBuilder()

const (
	ComingCollectionName = "Coming's Collection"
	ComingTokenName      = "Coming's Token"
)

func TestCreateCollection(t *testing.T) {
	payload, err := nftBuilder.CreateCollection(
		ComingCollectionName,
		"This is a collection",
		"https://www.comingchat.com",
		10)
	require.Nil(t, err)
	t.Log(payload)
}

func TestCreateToken(t *testing.T) {
	payload, err := nftBuilder.CreateToken(
		ComingCollectionName,
		ComingTokenName,
		"This is a token",
		"https://aptos.dev/img/nyan.jpeg",
		1,
		0,
		NFTRoyalty{},
		nil)
	require.Nil(t, err)
	t.Log(payload)
}

func TestOfferToken(t *testing.T) {
	receiver, _ := txnBuilder.NewAccountAddressFromHex("0x06070809")
	creator, _ := txnBuilder.NewAccountAddressFromHex("0x01020304")
	payload, err := nftBuilder.OfferToken(
		*receiver,
		*creator,
		ComingCollectionName,
		ComingTokenName,
		1,
		0,
	)
	require.Nil(t, err)
	t.Log(payload)
}

func TestClaimToken(t *testing.T) {
	sender, _ := txnBuilder.NewAccountAddressFromHex("0x06070809")
	creator, _ := txnBuilder.NewAccountAddressFromHex("0x01020304")
	payload, err := nftBuilder.ClaimToken(
		*sender,
		*creator,
		ComingCollectionName,
		ComingTokenName,
		1,
	)
	require.Nil(t, err)
	t.Log(payload)
}
