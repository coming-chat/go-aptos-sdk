package aptosclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"
	"os/exec"
	"strconv"
	"strings"
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
	ReceiverAddress = "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
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

	client, err := Dial(context.Background(), RestUrl)
	require.Nil(t, err)
	balance, err := client.AptosBalanceOf(address)
	require.Nil(t, err)
	t.Log(balance)
}

func TestTransferBCS(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	client, err := Dial(context.Background(), RestUrl)
	require.Nil(t, err)

	ledgerInfo, err := client.LedgerInfo()
	require.Nil(t, err)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	require.Nil(t, err)

	txn, err := generateTransactionBcs(accountData, ledgerInfo, account.AuthKey, toAddress, amount)
	require.Nil(t, err)

	signedTxn, err := txBuilder.GenerateBCSTransaction(account, txn)
	require.Nil(t, err)

	if !txnSubmitableForTest(t) {
		return
	}
	newTxn, err := client.SubmitSignedBCSTransaction(signedTxn)
	require.Nil(t, err)

	t.Logf("submited tx hash = %v", newTxn.Hash)
}

func TestBCSEncoder(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	client, err := Dial(context.Background(), RestUrl)
	require.Nil(t, err)

	ledgerInfo, err := client.LedgerInfo()
	require.Nil(t, err)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	require.Nil(t, err)

	// txn json
	txnJson, err := generateTransactionJson(accountData, ledgerInfo, account, toAddress, amount)
	require.Nil(t, err)

	signingMessageFromJson, err := client.CreateTransactionSigningMessage(txnJson)
	require.Nil(t, err)

	// txn bcs
	txnBcs, err := generateTransactionBcs(accountData, ledgerInfo, account.AuthKey, toAddress, amount)
	require.Nil(t, err)

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

	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	client, err := Dial(context.Background(), RestUrl)
	require.Nil(t, err)

	ledgerInfo, err := client.LedgerInfo()
	require.Nil(t, err)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	require.Nil(t, err)

	txn, err := generateTransactionBcs(accountData, ledgerInfo, account.AuthKey, toAddress, amount)
	require.Nil(t, err)

	signedTxn, err := txBuilder.GenerateBCSSimulation(account.PublicKey, txn)
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

func generateTransactionBcs(
	data *aptostypes.AccountCoreData,
	info *aptostypes.LedgerInfo,
	fromAuthkey [32]byte,
	to string, amount uint64) (txn *txBuilder.RawTransaction, err error) {

	moduleName, err := txBuilder.NewModuleIdFromString("0x1::coin")
	if err != nil {
		return
	}
	token, err := txBuilder.NewTypeTagStructFromString("0x1::aptos_coin::AptosCoin")
	if err != nil {
		return
	}
	toAddr, err := txBuilder.NewAccountAddressFromHex(to)
	if err != nil {
		return
	}
	toAmountBytes := txBuilder.BCSSerializeBasicValue(amount)
	payload := txBuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "transfer",
		TyArgs:       []txBuilder.TypeTag{*token},
		Args: [][]byte{
			toAddr[:], toAmountBytes,
		},
	}
	txn = &txBuilder.RawTransaction{
		Sender:                  fromAuthkey,
		SequenceNumber:          data.SequenceNumber,
		Payload:                 payload,
		MaxGasAmount:            2000,
		GasUnitPrice:            1,
		ExpirationTimestampSecs: info.LedgerTimestamp + 600,
		ChainId:                 uint8(info.ChainId),
	}
	return
}

func generateTransactionJson(
	data *aptostypes.AccountCoreData,
	info *aptostypes.LedgerInfo,
	from *aptosaccount.Account,
	to string, amount uint64) (txn *aptostypes.Transaction, err error) {

	amountString := strconv.FormatUint(amount, 10)
	payload := &aptostypes.Payload{
		Type:          aptostypes.EntryFunctionPayload,
		Function:      "0x1::coin::transfer",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		// Function:      "0x1::account::transfer",
		// TypeArguments: []string{},
		Arguments: []interface{}{
			to, amountString,
		},
	}
	fromAddress := "0x" + hex.EncodeToString(from.AuthKey[:])
	txn = &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          data.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            1,
		Payload:                 payload,
		ExpirationTimestampSecs: info.LedgerTimestamp + 600,
	}
	return
}

func txnSubmitableForTest(t *testing.T) bool {
	out, _ := exec.Command("whoami").Output()
	user := strings.TrimSpace(string(out))
	switch user {
	case "gg":
		return true
	default:
		t.Log("Non-specified machines, stop sending transactions after signing: ", user)
		return false
	}
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

	client, err := Dial(context.Background(), RestUrl)
	require.Nil(t, err)
	ensureBalanceGreatherThan(t, client, msPubkey.Address(), 2000)

	ledgerInfo, err := client.LedgerInfo()
	require.Nil(t, err)
	accountData, err := client.GetAccount(msPubkey.Address())
	require.Nil(t, err)
	txn, err := generateTransactionBcs(accountData, ledgerInfo, msPubkey.AuthenticationKey(), ReceiverAddress, 800)
	require.Nil(t, err)

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
	client, err := Dial(context.Background(), RestUrl)
	require.Nil(t, err)

	price, err := client.EstimateGasPrice()
	require.Nil(t, err)

	t.Log(price)
}
