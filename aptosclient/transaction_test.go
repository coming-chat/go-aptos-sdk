package aptosclient

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/go-aptos/aptostypes"
	txBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/lcs"
	"github.com/stretchr/testify/require"
)

const (
	Mnemonic        = "crack coil okay hotel glue embark all employ east impact stomach cigar"
	MnemonicAddress = "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
	ReceiverAddress = "0x6ed6f83f1891e02c00c58bf8172e3311c982b1c4fbb1be2d85a55562d4085fb1"
)

func TestFaucet(t *testing.T) {
	// address := ReceiverAddress
	address := MnemonicAddress
	hashs, err := FaucetFundAccount(address, 1000, "")
	require.Nil(t, err)
	t.Log(hashs)
}

func TestAccountBalance(t *testing.T) {
	address := MnemonicAddress

	client := Client(t, DevnetRestUrl)
	balance, err := client.AptosBalanceOf(address)
	require.Nil(t, err)
	t.Log(balance)
}

func TestTransferBCS(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	client := Client(t, DevnetRestUrl)
	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	params := transferParams{}
	params.transferFrom(t, account.AuthKey, client)
	params.transferTo(toAddress, amount)
	bcsTxn := params.generateTransactionBcs(t)

	signedTxn, err := txBuilder.GenerateBCSTransaction(account, bcsTxn)
	require.Nil(t, err)

	if !txnSubmitableForTest(t) {
		return
	}
	newTxn, err := client.SubmitSignedBCSTransaction(signedTxn)
	require.Nil(t, err)

	t.Logf("submited tx hash = %v", newTxn.Hash)
}

func TestTransferBatch(t *testing.T) {
	toAddress := "0xd408eeba43ada8ec52f1566518443950826450c16a079dfc11c2991b0f052f60"

	client := Client(t, DevnetRestUrl)
	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	txns := make([]*txBuilder.RawTransaction, 0)
	// 1
	params1 := transferParams{}
	params1.transferFrom(t, account.AuthKey, client)
	params1.transferTo(toAddress, 100)
	txns = append(txns, params1.generateTransactionBcs(t))

	// 2
	params2 := transferParams{}
	params2.transferFrom(t, account.AuthKey, client)
	params2.transferTo(toAddress, 200)
	rawTx2 := params2.generateTransactionBcs(t)
	rawTx2.SequenceNumber += 1
	txns = append(txns, rawTx2)

	// 3
	params3 := transferParams{}
	params3.transferFrom(t, account.AuthKey, client)
	params3.transferTo(toAddress, 300)
	rawTx3 := params3.generateTransactionBcs(t)
	rawTx3.SequenceNumber += 2
	txns = append(txns, rawTx3)
	// 4
	params4 := transferParams{}
	params4.transferFrom(t, account.AuthKey, client)
	params4.transferTo(toAddress, 400)
	rawTx4 := params4.generateTransactionBcs(t)
	rawTx4.SequenceNumber += 3
	txns = append(txns, rawTx4)
	// 5
	params5 := transferParams{}
	params5.transferFrom(t, account.AuthKey, client)
	params5.transferTo(toAddress, 500)
	rawTx5 := params5.generateTransactionBcs(t)
	rawTx5.SequenceNumber += 4
	txns = append(txns, rawTx5)

	signedTxn, err := txBuilder.GenerateBCSTransactionBatch(account, txns)
	require.Nil(t, err)

	res, err := client.SubmitSignedBCSTransactionBatch(signedTxn)
	require.Nil(t, err)

	t.Logf("submited tx hash = %v", res.TransactionFailures)
}

func TestBCSEncoder(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	client := Client(t, DevnetRestUrl)
	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	params := transferParams{}
	params.transferFrom(t, account.AuthKey, client)
	params.transferTo(toAddress, amount)

	// txn json
	txnJson := params.generateTransactionJson()
	signingMessageFromJson, err := client.CreateTransactionSigningMessage(txnJson)
	require.Nil(t, err)

	// txn bcs
	txnBcs := params.generateTransactionBcs(t)
	signingMessageFromBcs, err := txnBcs.GetSigningMessage()
	require.Nil(t, err)

	// compare bcs encoded results between remote server and local.
	hexStringBcs := hex.EncodeToString(signingMessageFromBcs)
	hexStringJson := hex.EncodeToString(signingMessageFromJson)
	require.Equal(t, hexStringBcs, hexStringJson)
}

func TestEstimateTransactionFeeBcs(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	client := Client(t, DevnetRestUrl)
	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	params := transferParams{}
	params.transferFrom(t, account.AuthKey, client)
	params.transferTo(toAddress, amount)

	bcsTxn := params.generateTransactionBcs(t)
	signedTxn, err := txBuilder.GenerateBCSSimulation(account.PublicKey, bcsTxn)
	require.Nil(t, err)
	newTxns, err := client.SimulateSignedBCSTransaction(signedTxn)
	require.Nil(t, err)

	if len(newTxns) == 0 {
		t.Fatal("simlated txn count empty")
	}
	firstTxn := newTxns[0]
	t.Logf("simlated tx hash = %v", firstTxn.Hash)
	t.Logf("gas price = %v, gas used = %v", firstTxn.GasUnitPrice, firstTxn.GasUsed)
}

func TestMultiSignTransfer(t *testing.T) {
	pri1 := [32]byte{1}
	pri2 := [32]byte{2}
	pri3 := [32]byte{3}
	account1 := aptosaccount.NewAccount(pri1[:])
	account2 := aptosaccount.NewAccount(pri2[:])
	account3 := aptosaccount.NewAccount(pri3[:])
	msPubkey, err := txBuilder.NewMultiEd25519PublicKey([][]byte{
		account1.PublicKey,
		account2.PublicKey,
		account3.PublicKey,
	}, 2)
	t.Logf("%x", msPubkey.ToBytes())
	t.Logf("%v", msPubkey.Address())

	client := Client(t, DevnetRestUrl)
	ensureBalanceGreatherThan(t, &client, msPubkey.Address(), 2000)

	params := transferParams{}
	params.transferFrom(t, msPubkey.AuthenticationKey(), client)
	params.transferTo(ReceiverAddress, 800)
	txn := params.generateTransactionBcs(t)

	// sign one by one
	signatures := [][]byte{}
	idxes := []uint8{}
	accountSigning := func(account *aptosaccount.Account, rawTxn *txBuilder.RawTransaction) {
		idx := indexOfPubkey(msPubkey, account.PublicKey)
		require.NotEqual(t, -1, idx, "the account not the member of the multi sign")
		signingMsg, err := rawTxn.GetSigningMessage()
		require.Nil(t, err)
		sign := account.Sign(signingMsg, "")

		signatures = append(signatures, sign)
		idxes = append(idxes, uint8(idx))
	}

	// Can be signed in any order: [1, 2], [3, 1], [3, 2], ...
	accountSigning(account3, txn)
	accountSigning(account1, txn)

	msSignature, err := txBuilder.NewMultiEd25519Signature(signatures, idxes)
	require.Nil(t, err)
	authenticator := txBuilder.TransactionAuthenticatorMultiEd25519{
		PublicKey: *msPubkey,
		Signature: *msSignature,
	}
	signedTxn := txBuilder.SignedTransaction{
		Transaction:   txn,
		Authenticator: authenticator,
	}
	signedTxnBytes, err := lcs.Marshal(signedTxn)
	require.Nil(t, err)

	// batch sign with builder
	// builder := txBuilder.TransactionBuilderMultiEd25519{
	// 	SigningFn: func(sm txBuilder.SigningMessage) txBuilder.MultiEd25519Signature {
	// 		sig1 := account1.Sign(sm, "")
	// 		sig3 := account3.Sign(sm, "")

	// 		signature, err := txBuilder.NewMultiEd25519Signature([][]byte{sig1, sig3}, []uint8{0, 2})
	// 		require.Nil(t, err)
	// 		return *signature
	// 	},
	// 	PublicKey: *msPubkey,
	// }
	// signedTxnBytes, err := builder.Sign(txn)
	// require.Nil(t, err)

	if !txnSubmitableForTest(t) {
		return
	}
	newTxn, err := client.SubmitSignedBCSTransaction(signedTxnBytes)
	require.Nil(t, err)

	t.Logf("multi sign transaction success: %v\n hash = %v", newTxn, newTxn.Hash)
}

func indexOfPubkey(msPubkey *txBuilder.MultiEd25519PublicKey, pubkey []byte) int {
	for idx, pub := range msPubkey.PublicKeys {
		if bytes.Compare(pubkey, pub.PublicKey) == 0 {
			return idx
		}
	}
	return -1
}

func ensureBalanceGreatherThan(t *testing.T, client *RestClient, address string, amount uint64) {
	balance, err := client.AptosBalanceOf(address)
	require.Nil(t, err)
	if balance.Cmp(big.NewInt(int64(amount))) < 0 {
		_, err = FaucetFundAccount(address, amount, "")
		require.Nil(t, err)
	}
}

func TestGasPrice(t *testing.T) {
	client := Client(t, DevnetRestUrl)
	price, err := client.EstimateGasPrice()
	require.Nil(t, err)
	t.Log(price)
}

type transferParams struct {
	Receiver string
	Amount   uint64

	SenderKey   [32]byte
	AccountData *aptostypes.AccountCoreData
	LedgerInfo  *aptostypes.LedgerInfo
	GasPrice    uint64
}

func (p *transferParams) transferFrom(t *testing.T, sender [32]byte, client RestClient) {
	p.SenderKey = sender
	address := "0x" + hex.EncodeToString(sender[:])

	var err error = nil
	p.LedgerInfo, err = client.LedgerInfo()
	require.Nil(t, err)
	p.AccountData, err = client.GetAccount(address)
	require.Nil(t, err)
	p.GasPrice, err = client.EstimateGasPrice()
	require.Nil(t, err)
}

func (p *transferParams) transferTo(receiver string, amount uint64) {
	p.Receiver = receiver
	p.Amount = amount
}

func (p *transferParams) generateTransactionBcs(t *testing.T) *txBuilder.RawTransaction {
	moduleName, err := txBuilder.NewModuleIdFromString("0x1::coin")
	require.Nil(t, err)
	token, err := txBuilder.NewTypeTagStructFromString("0x1::aptos_coin::AptosCoin")
	require.Nil(t, err)
	toAddr, err := txBuilder.NewAccountAddressFromHex(p.Receiver)
	require.Nil(t, err)
	toAmountBytes := txBuilder.BCSSerializeBasicValue(p.Amount)
	payload := txBuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "transfer",
		TyArgs:       []txBuilder.TypeTag{*token},
		Args: [][]byte{
			toAddr[:], toAmountBytes,
		},
	}
	return &txBuilder.RawTransaction{
		Sender:                  p.SenderKey,
		SequenceNumber:          p.AccountData.SequenceNumber,
		Payload:                 payload,
		MaxGasAmount:            2000,
		GasUnitPrice:            p.GasPrice,
		ExpirationTimestampSecs: p.LedgerInfo.LedgerTimestamp + 600,
		ChainId:                 uint8(p.LedgerInfo.ChainId),
	}
}

func (p *transferParams) generateTransactionJson() *aptostypes.Transaction {
	amountString := strconv.FormatUint(p.Amount, 10)
	payload := &aptostypes.Payload{
		Type:          aptostypes.EntryFunctionPayload,
		Function:      "0x1::coin::transfer",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		// Function:      "0x1::aptos_account::transfer",
		// TypeArguments: []string{},
		Arguments: []interface{}{
			p.Receiver, amountString,
		},
	}
	fromAddress := "0x" + hex.EncodeToString(p.SenderKey[:])
	return &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          p.AccountData.SequenceNumber,
		Payload:                 payload,
		MaxGasAmount:            2000,
		GasUnitPrice:            p.GasPrice,
		ExpirationTimestampSecs: p.LedgerInfo.LedgerTimestamp + 600,
	}
}
