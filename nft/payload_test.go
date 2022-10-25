package nft

import (
	"context"
	"encoding/hex"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
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

var (
	PriMartian1 = strings.TrimPrefix(os.Getenv("PriMartian1"), "0x")
	PriMartian2 = strings.TrimPrefix(os.Getenv("PriMartian2"), "0x")
	PriPetra1   = strings.TrimPrefix(os.Getenv("PriPetra1"), "0x")
)

func IsSupportedMachine(t *testing.T) bool {
	ok := PriMartian1 != "" && PriMartian2 != "" && PriPetra1 != ""
	if !ok {
		t.Log("The tested machine is not supported, env private key not found.")
	}
	return ok
}

func TestNFTTransfer(t *testing.T) {
	if !IsSupportedMachine(t) {
		return
	}

	// pri1, err := hex.DecodeString(PriMartian1)
	pri1, err := hex.DecodeString(PriMartian2)
	require.Nil(t, err)
	acc1 := aptosaccount.NewAccount(pri1)
	t.Log("address1:", hex.EncodeToString(acc1.AuthKey[:]))

	pri2, err := hex.DecodeString(PriPetra1)
	require.Nil(t, err)
	acc2 := aptosaccount.NewAccount(pri2)
	t.Log("address2:", hex.EncodeToString(acc2.AuthKey[:]))

	// 测试目标：
	// 让 account1 发 offer nft 给 account2
	// 然后 account2 发 claim nft
	client, err := aptosclient.Dial(context.Background(), TestnetRestUrl)
	require.Nil(t, err)
	tokenClient := NewTokenClient(client)
	nfts, err := tokenClient.GetAllTokenForAccount(acc1.AuthKey)
	require.Nil(t, err)
	if len(nfts) == 0 {
		t.Log("break off: because account1 have not nft.")
		return
	}
	tokenId := nfts[0].TokenId
	t.Log("=== transferring token info: ", tokenId)

	// account1 offer nft to account2
	offerTxn := OfferToken(t, *client, *acc1, *acc2, *tokenId)
	t.Log("offer hash:", offerTxn.Hash)

	time.Sleep(time.Second * 10)
	detail, err := client.GetTransactionByHash(offerTxn.Hash)
	require.Nil(t, err)
	if !detail.Success {
		t.Log("offer token failed, stop claim action. ", detail.VmStatus)
		return
	}

	// account2 claim nft from account1
	claimTxn := ClaimToken(t, *client, *acc1, *acc2, *tokenId)
	t.Log("claim hash:", claimTxn.Hash)
}

func OfferToken(t *testing.T, client aptosclient.RestClient, sender, receiver aptosaccount.Account, tokenId TokenDataId) *aptostypes.Transaction {
	creator, err := txnBuilder.NewAccountAddressFromHex(tokenId.Creator)
	require.Nil(t, err)
	offerPayload, err := nftBuilder.OfferToken(receiver.AuthKey, *creator, tokenId.Collection, tokenId.Name, 1, 0)
	require.Nil(t, err)
	return SignAndSendPayload(t, client, sender, offerPayload)
}
func ClaimToken(t *testing.T, client aptosclient.RestClient, sender, receiver aptosaccount.Account, tokenId TokenDataId) *aptostypes.Transaction {
	creator, err := txnBuilder.NewAccountAddressFromHex(tokenId.Creator)
	require.Nil(t, err)
	claimPayload, err := nftBuilder.ClaimToken(sender.AuthKey, *creator, tokenId.Collection, tokenId.Name, 0)
	require.Nil(t, err)
	return SignAndSendPayload(t, client, receiver, claimPayload)
}

func SignAndSendPayload(t *testing.T, client aptosclient.RestClient, account aptosaccount.Account, payload txnBuilder.TransactionPayload) *aptostypes.Transaction {
	senderAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	ledgerInfo, err := client.LedgerInfo()
	require.Nil(t, err)
	accountData, err := client.GetAccount(senderAddress)
	require.Nil(t, err)
	gasPrice, err := client.EstimateGasPrice()
	require.Nil(t, err)

	txn := txnBuilder.RawTransaction{
		Sender:                  account.AuthKey,
		SequenceNumber:          accountData.SequenceNumber,
		Payload:                 payload,
		MaxGasAmount:            2000,
		GasUnitPrice:            gasPrice,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
		ChainId:                 uint8(ledgerInfo.ChainId),
	}
	signedTxn, err := txnBuilder.GenerateBCSTransaction(&account, &txn)
	require.Nil(t, err)
	newTxn, err := client.SubmitSignedBCSTransaction(signedTxn)
	require.Nil(t, err)
	return newTxn
}
