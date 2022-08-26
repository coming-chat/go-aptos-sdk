package aptosclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/go-aptos/aptostypes"
	txBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/the729/lcs"
)

const (
	Mnemonic        = "crack coil okay hotel glue embark all employ east impact stomach cigar"
	ReceiverAddress = "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
)

func TestTransferBCS(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	checkError(t, err)

	client, err := Dial(context.Background(), RestUrl)
	checkError(t, err)

	ledgerInfo, err := client.LedgerInfo()
	checkError(t, err)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	checkError(t, err)

	txn, err := generateTransactionBcs(accountData, ledgerInfo, account, toAddress, amount)
	checkError(t, err)

	signedTxn, err := txBuilder.GenerateBCSTransaction(account, txn)
	checkError(t, err)

	if !txnSubmitableForTest(t) {
		return
	}
	newTxn, err := client.SubmitSignedBCSTransaction(signedTxn)
	checkError(t, err)

	t.Logf("submited tx hash = %v", newTxn.Hash)
}

func TestTransferJson(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	checkError(t, err)

	client, err := Dial(context.Background(), RestUrl)
	checkError(t, err)

	ledgerInfo, err := client.LedgerInfo()
	checkError(t, err)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	checkError(t, err)

	txn, err := generateTransactionJson(accountData, ledgerInfo, account, toAddress, amount)
	checkError(t, err)

	signingMessage, err := client.CreateTransactionSigningMessage(txn)
	checkError(t, err)

	signatureData := account.Sign(signingMessage, "")
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	txn.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}

	if !txnSubmitableForTest(t) {
		return
	}
	newTxn, err := client.SubmitTransaction(txn)
	checkError(t, err)

	t.Logf("submited tx hash = %v", newTxn.Hash)
}

func TestBCSEncoder(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	checkError(t, err)

	client, err := Dial(context.Background(), RestUrl)
	checkError(t, err)

	ledgerInfo, err := client.LedgerInfo()
	checkError(t, err)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	checkError(t, err)

	// txn json
	txnJson, err := generateTransactionJson(accountData, ledgerInfo, account, toAddress, amount)
	checkError(t, err)

	signingMessageFromJson, err := client.CreateTransactionSigningMessage(txnJson)
	checkError(t, err)

	// txn bcs
	txnBcs, err := generateTransactionBcs(accountData, ledgerInfo, account, toAddress, amount)
	checkError(t, err)

	signingMessageFromBcs, err := txnBcs.GetSigningMessage()
	checkError(t, err)

	// compare bcs encoded results between remote server and local.
	if bytes.Compare(signingMessageFromJson, signingMessageFromBcs) == 0 {
		t.Logf("signingMessage = %x", signingMessageFromBcs)
		t.Log("generate bcs passed")
	} else {
		t.Logf("json = %x", signingMessageFromJson)
		t.Logf("bcs  = %x", signingMessageFromBcs)
		t.Fatal("generate bcs failed")
	}
}

func TestEstimateTransactionFeeBcs(t *testing.T) {
	toAddress := ReceiverAddress
	amount := uint64(100)

	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	checkError(t, err)

	client, err := Dial(context.Background(), RestUrl)
	checkError(t, err)

	ledgerInfo, err := client.LedgerInfo()
	checkError(t, err)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	checkError(t, err)

	txn, err := generateTransactionBcs(accountData, ledgerInfo, account, toAddress, amount)
	checkError(t, err)

	signedTxn, err := txBuilder.GenerateBCSSimulation(account, txn)
	checkError(t, err)

	newTxns, err := client.SimulateSignedBCSTransaction(signedTxn)
	checkError(t, err)

	if len(newTxns) == 0 {
		t.Fatal("simlated txn count empty")
	}
	firstTxn := newTxns[0]
	t.Logf("simlated tx hash = %v", firstTxn.Hash)
	t.Logf("gas price = %v, gas used = %v", firstTxn.GasUnitPrice, firstTxn.GasUsed)
}

func TestFaucet(t *testing.T) {
	address := ReceiverAddress
	hashs, err := FaucetFundAccount(address, 1000, "")
	checkError(t, err)
	t.Log(hashs)
}

func TestAccountBalance(t *testing.T) {
	address := ReceiverAddress

	client, err := Dial(context.Background(), RestUrl)
	checkError(t, err)
	balance, err := client.BalanceOf(address)
	checkError(t, err)
	t.Log(balance)
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func generateTransactionBcs(
	data *aptostypes.AccountCoreData,
	info *aptostypes.LedgerInfo,
	from *aptosaccount.Account,
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
	toAmountBytes, _ := lcs.Marshal(amount)
	payload := txBuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "transfer",
		TyArgs:       []txBuilder.TypeTag{*token},
		Args: [][]byte{
			toAddr[:], toAmountBytes,
		},
	}
	txn = &txBuilder.RawTransaction{
		Sender:                  from.AuthKey,
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
