package transactionbuilder

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/stretchr/testify/assert"
)

const (
	Mnemonic        = "crack coil okay hotel glue embark all employ east impact stomach cigar"
	ReceiverAddress = "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
	RestUrl         = "https://fullnode.devnet.aptoslabs.com/"
)

func TestTransactionBuilderABI(t *testing.T) {
	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	assert.Nil(t, err)
	client, err := aptosclient.Dial(context.Background(), RestUrl)
	assert.Nil(t, err)

	ledgerInfo, err := client.LedgerInfo()
	assert.Nil(t, err)
	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	assert.Nil(t, err)

	functionName := "0xb39c45e31d1429218aeb3590e2a046edae9303fbbc3ef6a065384569cfd81881::red_packet::create"

	// build transaction with json
	payloadJson := &aptostypes.Payload{
		Type:          aptostypes.EntryFunctionPayload,
		Function:      functionName,
		TypeArguments: []string{},
		Arguments: []interface{}{
			"5", "1000000",
		},
	}
	txnJson := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            1,
		Payload:                 payloadJson,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
	}

	signingMessageJson, err := client.CreateTransactionSigningMessage(txnJson)
	assert.Nil(t, err)
	t.Logf("%x: signing message from json", signingMessageJson)

	// build transaction with abi
	redpacketABI := LoadRedPacketABI(t)
	payloadABI, err := redpacketABI.BuildTransactionPayload(
		functionName,
		[]string{},
		[]any{
			uint64(5), uint64(1e6),
		},
	)
	assert.Nil(t, err)
	txnABI := &RawTransaction{
		Sender:                  account.AuthKey,
		SequenceNumber:          accountData.SequenceNumber,
		Payload:                 payloadABI,
		MaxGasAmount:            2000,
		GasUnitPrice:            1,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
		ChainId:                 uint8(ledgerInfo.ChainId),
	}

	signingMessageABI, err := txnABI.GetSigningMessage()
	assert.Nil(t, err)
	t.Logf("%x: signing message from abi ", signingMessageABI)

	// compare
	assert.Equal(t, SigningMessage(signingMessageJson), signingMessageABI, "The signingMessage from abi must be the same")
}

func TestTransactionBuilderABIOpen(t *testing.T) {
	account, err := aptosaccount.NewAccountWithMnemonic(Mnemonic)
	assert.Nil(t, err)
	client, err := aptosclient.Dial(context.Background(), RestUrl)
	assert.Nil(t, err)

	ledgerInfo, err := client.LedgerInfo()
	assert.Nil(t, err)
	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	accountData, err := client.GetAccount(fromAddress)
	assert.Nil(t, err)

	functionName := "0xb39c45e31d1429218aeb3590e2a046edae9303fbbc3ef6a065384569cfd81881::red_packet::open"

	// build transaction with json
	payloadJson := &aptostypes.Payload{
		Type:          aptostypes.EntryFunctionPayload,
		Function:      functionName,
		TypeArguments: []string{},
		Arguments: []interface{}{
			"5", []string{AccountAddress{0x1}.ToShortString(), "0x22"}, []string{"100", "200"},
		},
	}
	txnJson := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            1,
		Payload:                 payloadJson,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
	}

	signingMessageJson, err := client.CreateTransactionSigningMessage(txnJson)
	assert.Nil(t, err)
	t.Logf("%x: signing message from json", signingMessageJson)

	// build transaction with abi
	redpacketABI := LoadRedPacketABI(t)
	payloadABI, err := redpacketABI.BuildTransactionPayload(
		functionName,
		[]string{},
		[]any{
			uint64(5), []interface{}{AccountAddress{0x1}, "0x22"}, []uint64{100, 200},
		},
	)
	assert.Nil(t, err)
	txnABI := &RawTransaction{
		Sender:                  account.AuthKey,
		SequenceNumber:          accountData.SequenceNumber,
		Payload:                 payloadABI,
		MaxGasAmount:            2000,
		GasUnitPrice:            1,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
		ChainId:                 uint8(ledgerInfo.ChainId),
	}

	signingMessageABI, err := txnABI.GetSigningMessage()
	assert.Nil(t, err)
	t.Logf("%x: signing message from abi ", signingMessageABI)

	// compare
	assert.Equal(t, SigningMessage(signingMessageJson), signingMessageABI, "The signingMessage from abi must be the same")
}

func LoadRedPacketABI(t *testing.T) *TransactionBuilderABI {
	abiHexStrings := []string{
		// create
		"0106637265617465b39c45e31d1429218aeb3590e2a046edae9303fbbc3ef6a065384569cfd818810a7265645f7061636b657400000205636f756e74020d746f74616c5f62616c616e636502",
		// open
		"01046f70656eb39c45e31d1429218aeb3590e2a046edae9303fbbc3ef6a065384569cfd818810a7265645f7061636b6574000003026964020e6c75636b795f6163636f756e747306040862616c616e6365730602",
		// close
		"0105636c6f7365b39c45e31d1429218aeb3590e2a046edae9303fbbc3ef6a065384569cfd818810a7265645f7061636b657400000102696402",
	}
	abiBytes := [][]byte{}
	for _, hexString := range abiHexStrings {
		bytes, err := hex.DecodeString(hexString)
		assert.Nil(t, err)
		abiBytes = append(abiBytes, bytes)
	}

	redPacketABI, err := NewTransactionBuilderABI(abiBytes)
	assert.Nil(t, err)

	return redPacketABI
}
