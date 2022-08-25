package transactionbuilder

import (
	"fmt"
)

const (
	ED25519_PUBLICKEY_LENGTH = 32
	ED25519_SIGNATURE_LENGTH = 64

	// MULTI_ED25519_SIGNATURE_BITMAP_LENGTH = 4
	// MAX_SIGNATURE_SUPPORTED = 32
)

// ------ PublicKey ------

type PublicKey interface{}

type Ed25519PublicKey struct {
	PublicKey []byte `lcs:"publicKey"`
}

func NewEd25519PublicKey(publicKey []byte) (*Ed25519PublicKey, error) {
	if len(publicKey) != ED25519_PUBLICKEY_LENGTH {
		return nil, fmt.Errorf(`Ed25519PublicKey length should be %d`, ED25519_PUBLICKEY_LENGTH)
	}
	return &Ed25519PublicKey{publicKey}, nil
}

// ------ Signature ------

type Signature interface{}

type Ed25519Signature struct {
	Signature []byte `lcs:"signature"`
}

func NewEd25519Signature(signature []byte) (*Ed25519Signature, error) {
	if len(signature) != ED25519_SIGNATURE_LENGTH {
		return nil, fmt.Errorf(`Ed25519Signature length should be %d`, ED25519_SIGNATURE_LENGTH)
	}
	return &Ed25519Signature{signature}, nil
}

// type MultiEd25519Signature struct {
// 	Bytes []byte `lcs:"bytes"`
// 	Signatures []Ed25519Signature
// 	Bitmap     []byte
// }
