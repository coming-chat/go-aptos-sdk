# Aptos Goland SDK

[![Documentation (master)](https://img.shields.io/badge/docs-master-59f)](https://github.com/coming-chat/wallet-SDK)
[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://github.com/aptos-labs/aptos-core/blob/main/LICENSE)

The Aptos Golang SDK for ComingChat.
We welcome other developers to participate in the development and testing of golang-sdk @aptos-labs！

## Install

```sh
go get github.com/coming-chat/go-aptos
```

## Usage

### Account

```go
import "github.com/coming-chat/go-aptos/aptosaccount"

// Import account with mnemonic
account, err := aptosaccount.NewAccountWithMnemonic(mnemonic)

// Import account with private key
privateKey, err := hex.DecodeString("4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b")
account, err := aptosaccount.NewAccount(privateKey)

// Get private key, public key, address
fmt.Printf("privateKey = %x\n", account.PrivateKey[:32])
fmt.Printf(" publicKey = %x\n", account.PublicKey)
fmt.Printf("   address = %x\n", account.AuthKey)

// Sign data
signedData := account.Sign(data, "")
```

### Transfer Aptos Coin (BCS)
> NOTE: The implementation of bcs is based on the [lcs](https://github.com/coming-chat/lcs) library

```go
import "github.com/coming-chat/go-aptos/aptosaccount"
import "github.com/coming-chat/go-aptos/aptosclient"
import "github.com/coming-chat/go-aptos/aptostypes"
import txBuilder "github.com/coming-chat/go-aptos/transaction_builder"

account, err := NewAccountWithMnemonic(mnemonic)
fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])

toAddress := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
amount := uint64(100)

// Initialize the client
restUrl := "https://fullnode.devnet.aptoslabs.com"
client, err := aptosclient.Dial(context.Background(), restUrl)

// Get Sender's account data and ledger info
accountData, err := client.GetAccount(fromAddress)
ledgerInfo, err := client.LedgerInfo()

// Build paylod
moduleName, err := txBuilder.NewModuleIdFromString("0x1::account")
toAddr, err := txBuilder.NewAccountAddressFromHex(toAddress)
toAmountBytes := txBuilder.BCSSerializeBasicValue(amount)

payload := txBuilder.TransactionPayloadEntryFunction{
	ModuleName:   *moduleName,
	FunctionName: "transfer",
	TyArgs:       []txBuilder.TypeTag{},
	Args: [][]byte{
		toAddr[:], toAmountBytes,
	},
}

// Build transaction
txn = &txBuilder.RawTransaction{
	Sender:                  from.AuthKey,
	SequenceNumber:          data.SequenceNumber,
	Payload:                 payload,
	MaxGasAmount:            2000,
	GasUnitPrice:            1,
	ExpirationTimestampSecs: info.LedgerTimestamp + 600,
	ChainId:                 uint8(info.ChainId),
}

// Sign raw transaction with account, and encode into data using BCS
signedTxn, err := txBuilder.GenerateBCSTransaction(account, txn)

// Submit transaction with BCS format.
newTxn, err := client.SubmitSignedBCSTransaction(signedTxn)
fmt.Printf("tx hash = %v\n", newTx.Hash)
```

### Transfer Aptos Coin (JSON)

```go
import "github.com/coming-chat/go-aptos/aptosaccount"
import "github.com/coming-chat/go-aptos/aptosclient"
import "github.com/coming-chat/go-aptos/aptostypes"

account, err := NewAccountWithMnemonic(mnemonic)
fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])

toAddress := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
amount := "100"

// Initialize the client
restUrl := "https://fullnode.devnet.aptoslabs.com"
client, err := aptosclient.Dial(context.Background(), restUrl)

// Get Sender's account data and ledger info
accountData, err := client.GetAccount(fromAddress)
ledgerInfo, err := client.LedgerInfo()

// Build paylod
payload := &aptostypes.Payload{
	Type: 		   "entry_function_payload",
	Function:      "0x1::coin::transfer",
	TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
	Arguments: []interface{}{
		toAddress, amount,
	},
}

// Build transaction
transaction := &aptostypes.Transaction{
	Sender:                  fromAddress,
	SequenceNumber:          accountData.SequenceNumber,
	MaxGasAmount:            2000,
	GasUnitPrice:            1,
	Payload:                 payload,
	ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
}

// Get signing message from remote server
// Note: Later we will implement the local use of BCS encoding to create signing messages
signingMessage, err := client.CreateTransactionSigningMessage(transaction)

// Sign message and complete transaction information
signatureData := account.Sign(signingMessage, "")
signatureHex := "0x" + hex.EncodeToString(signatureData)
publicKey := "0x" + hex.EncodeToString(account.PublicKey)
transaction.Signature = &aptostypes.Signature{
	Type:      "ed25519_signature",
	PublicKey: publicKey,
	Signature: signatureHex,
}

// Submit transaction
newTx, err := client.SubmitTransaction(transaction)
fmt.Printf("tx hash = %v\n", newTx.Hash)
```



## TODO

- [x] Locally implement BCS encoding of coin transfer transaction data
- [ ] Support Move action BCS
- [ ] Support Multi sign
