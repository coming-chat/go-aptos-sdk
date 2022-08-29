package transactionbuilder

import "github.com/coming-chat/lcs"

func init() {
	lcs.RegisterEnum(
		(*TransactionAuthenticator)(nil),

		TransactionAuthenticatorEd25519{},
		// TransactionAuthenticatorMultiEd25519{},
		// TransactionAuthenticatorMultiAgent{},
	)

	lcs.RegisterEnum(
		(*AccountAuthenticator)(nil),

		AccountAuthenticatorEd25519{},
		// AccountAuthenticatorMultiEd25519{},
	)
}

// ------ TransactionAuthenticator ------

type TransactionAuthenticator interface{}

type TransactionAuthenticatorEd25519 struct {
	PublicKey Ed25519PublicKey `lcs:"publicKey"`
	Signature Ed25519Signature `lcs:"signature"`
}

// type TransactionAuthenticatorMultiEd25519 struct{}
// type TransactionAuthenticatorMultiAgent struct{}

// ------ AccountAuthenticator ------

type AccountAuthenticator interface{}

type AccountAuthenticatorEd25519 struct {
	PublicKey Ed25519PublicKey `lcs:"publicKey"`
	Signature Ed25519Signature `lcs:"signature"`
}

// type AccountAuthenticatorMultiEd25519 struct{}
