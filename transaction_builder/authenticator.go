package transactionbuilder

import "github.com/coming-chat/lcs"

func init() {
	lcs.RegisterEnum(
		(*TransactionAuthenticator)(nil),

		TransactionAuthenticatorEd25519{},
		TransactionAuthenticatorMultiEd25519{},
		TransactionAuthenticatorMultiAgent{},
	)

	lcs.RegisterEnum(
		(*AccountAuthenticator)(nil),

		AccountAuthenticatorEd25519{},
		AccountAuthenticatorMultiEd25519{},
	)
}

// ------ TransactionAuthenticator ------

type TransactionAuthenticator interface{}

type TransactionAuthenticatorEd25519 struct {
	PublicKey Ed25519PublicKey `lcs:"publicKey"`
	Signature Ed25519Signature `lcs:"signature"`
}

type TransactionAuthenticatorMultiEd25519 struct {
	PublicKey MultiEd25519PublicKey `lcs:"publicKey"`
	Signature MultiEd25519Signature `lcs:"signature"`
}

type TransactionAuthenticatorMultiAgent struct {
	Sender                   AccountAuthenticator   `lcs:"sender"`
	SecondarySignerAddresses []AccountAddress       `lcs:"secondary_signer_addresses"`
	SecondarySigners         []AccountAuthenticator `lcs:"secondary_signers"`
}

// ------ AccountAuthenticator ------

type AccountAuthenticator interface{}

type AccountAuthenticatorEd25519 struct {
	PublicKey Ed25519PublicKey `lcs:"publicKey"`
	Signature Ed25519Signature `lcs:"signature"`
}

type AccountAuthenticatorMultiEd25519 struct {
	PublicKey MultiEd25519PublicKey `lcs:"publicKey"`
	Signature MultiEd25519Signature `lcs:"signature"`
}
