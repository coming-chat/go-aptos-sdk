package transactionbuilder

import (
	"errors"

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

// ------ TransactionBuilder ------

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
