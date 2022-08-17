package aptosaccount

import (
	"context"
	"encoding/hex"
	"os/exec"
	"testing"

	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
)

const mnemonic = "dragon setup knee couch team journey genre barely nurse twelve blame toe"

func TestAccountSign(t *testing.T) {
	account, err := NewAccountWithMnemonic(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte{0x1}
	salt := "APTOS::RawTransaction"
	signature := account.Sign(data, salt)

	t.Logf("%x", signature)
}

func TestTransfer(t *testing.T) {
	account, err := NewAccountWithMnemonic(mnemonic)
	if err != nil {
		t.Fatal(err)
	}
	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])

	toAddress := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
	amount := "100"

	client, err := aptosclient.Dial(context.Background(), "https://fullnode.devnet.aptoslabs.com")
	if err != nil {
		t.Fatal(err)
	}

	accountData, err := client.GetAccount(fromAddress)
	if err != nil {
		t.Fatal(err)
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		t.Fatal(err)
	}

	payload := &aptostypes.Payload{
		Type: "script_function_payload",
		// Function:      "0x1::coin::transfer",
		// TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		Function:      "0x1::account::transfer",
		TypeArguments: []string{},
		Arguments: []interface{}{
			toAddress, amount,
		},
	}

	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            1,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
	}

	signingMessage, err := client.CreateTransactionSigningMessage(transaction)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signingMessage = %x", signingMessage)

	// const RAW_TRANSACTION_SALT = "APTOS::RawTransaction"
	signatureData := account.Sign(signingMessage, "")
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}
	t.Logf("signature = %x", signatureData)

	out, _ := exec.Command("whoami").Output()
	user := string(out)
	switch user {
	case "gg":
		break
	default:
		t.Log("Non-specified machines, stop sending transactions after signing: ", user)
		return
	}

	newTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		t.Fatal(err)
	}

	println("tx hash = ", newTx.Hash)
}

func TestAccountBalance(t *testing.T) {
	address := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"

	client, err := aptosclient.Dial(context.Background(), "https://fullnode.devnet.aptoslabs.com")
	if err != nil {
		t.Fatal(err)
	}

	balance, err := client.BalanceOf(address)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}

func TestFaucet(t *testing.T) {
	address := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
	hashs, err := aptosclient.FaucetFundAccount(address, 20, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashs)
}
