package transactionbuilder

import (
	"errors"
	"fmt"

	"github.com/coming-chat/lcs"
	"golang.org/x/crypto/sha3"
)

const (
	RAW_TRANSACTION_SALT           = "APTOS::RawTransaction"
	RAW_TRANSACTION_WITH_DATA_SALT = "APTOS::RawTransactionWithData"
)

type SigningMessage []byte

type Signable interface {
	GetSigningMessage() (SigningMessage, error)
}

func (t *RawTransaction) GetSigningMessage() (SigningMessage, error) {
	prefixBytes := sha3.Sum256([]byte(RAW_TRANSACTION_SALT))
	msg, err := lcs.Marshal(t)
	if err != nil {
		return nil, err
	}
	return append(prefixBytes[:], msg...), nil
}

func (t *MultiAgentRawTransaction) GetSigningMessage() (SigningMessage, error) {
	prefixBytes := sha3.Sum256([]byte(RAW_TRANSACTION_WITH_DATA_SALT))
	msg, err := lcs.Marshal(t)
	if err != nil {
		return nil, err
	}
	return append(prefixBytes[:], msg...), nil
}

// ------ TransactionBuilderEd25519 ------

type SigningFunctionEd25519 func(SigningMessage) []byte
type TransactionBuilderEd25519 struct {
	SigningFn SigningFunctionEd25519
	PublicKey []byte
}

func NewTransactionBuilderEd25519(signingFn SigningFunctionEd25519, publicKey []byte) *TransactionBuilderEd25519 {
	return &TransactionBuilderEd25519{signingFn, publicKey}
}

func (b *TransactionBuilderEd25519) Sign(rawTxn *RawTransaction) (data []byte, err error) {
	if b.SigningFn == nil {
		return nil, errors.New("Signing failed: you must specify a signing function")
	}
	signingMessage, err := rawTxn.GetSigningMessage()
	if err != nil {
		return
	}
	signatureBytes := b.SigningFn(signingMessage)
	if err != nil {
		return
	}
	publickey, err := NewEd25519PublicKey(b.PublicKey)
	if err != nil {
		return
	}
	signature, err := NewEd25519Signature(signatureBytes)
	if err != nil {
		return
	}
	authenticator := TransactionAuthenticatorEd25519{
		PublicKey: *publickey,
		Signature: *signature,
	}
	signedTxn := SignedTransaction{
		Transaction:   rawTxn,
		Authenticator: authenticator,
	}

	data, err = lcs.Marshal(signedTxn)
	return data, err
}

// ------ TransactionBuilderMultiEd25519 ------

type SigningFunctionMultiEd25519 func(SigningMessage) MultiEd25519Signature
type TransactionBuilderMultiEd25519 struct {
	SigningFn SigningFunctionMultiEd25519
	PublicKey MultiEd25519PublicKey
}

func (b *TransactionBuilderMultiEd25519) Sign(rawTxn *RawTransaction) (data []byte, err error) {
	if b.SigningFn == nil {
		return nil, errors.New("Signing failed: you must specify a signing function")
	}
	signingMessage, err := rawTxn.GetSigningMessage()
	if err != nil {
		return
	}
	signature := b.SigningFn(signingMessage)

	authenticator := TransactionAuthenticatorMultiEd25519{
		PublicKey: b.PublicKey,
		Signature: signature,
	}
	signedTxn := SignedTransaction{
		Transaction:   rawTxn,
		Authenticator: authenticator,
	}

	data, err = lcs.Marshal(signedTxn)
	return data, err
}

// ------ TransactionBuilderABI ------

type ABIBuilderConfig struct {
	Sender         AccountAddress
	SequenceNumber uint64
	GasUnitPrice   uint64
	MaxGasAmount   uint64
	ExpSecFromNow  uint64
	ChainId        uint8
}

type TransactionBuilderABI struct {
	ABIMap         map[string]ScriptABI
	BuildereConfig ABIBuilderConfig
}

func NewTransactionBuilderABI(abis [][]byte, config *ABIBuilderConfig) (*TransactionBuilderABI, error) {
	abiMap := make(map[string]ScriptABI)
	for _, bytes := range abis {
		var abi ScriptABI
		err := lcs.Unmarshal(bytes, &abi)
		if err != nil {
			return nil, err
		}

		k := ""
		if funcABI, ok := abi.(EntryFunctionABI); ok {
			module := funcABI.ModuleName
			k = fmt.Sprintf("%v::%v::%v", module.Address.ToShortString(), module.Name, funcABI.Name)
		} else {
			funcABI := abi.(TransactionScriptABI)
			k = funcABI.Name
		}

		if abiMap[k] != nil {
			return nil, errors.New("Found conflicting ABI interfaces")
		}
		abiMap[k] = abi
	}

	bc := ABIBuilderConfig{
		GasUnitPrice:  1,
		MaxGasAmount:  2000,
		ExpSecFromNow: 20,
	}
	if config != nil {
		bc.Sender = config.Sender
		bc.SequenceNumber = config.SequenceNumber
		bc.ChainId = config.ChainId
		if config.GasUnitPrice > 0 {
			bc.GasUnitPrice = config.GasUnitPrice
		}
		if config.MaxGasAmount > 0 {
			bc.MaxGasAmount = config.MaxGasAmount
		}
		if config.ExpSecFromNow > 0 {
			bc.ExpSecFromNow = config.ExpSecFromNow
		}
	}

	return &TransactionBuilderABI{
		ABIMap:         abiMap,
		BuildereConfig: bc,
	}, nil
}
