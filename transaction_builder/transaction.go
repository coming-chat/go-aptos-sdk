package transactionbuilder

import (
	"crypto/ed25519"
	"errors"
	"strings"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/lcs"
)

func init() {
	lcs.RegisterEnum(
		(*TransactionPayload)(nil),

		TransactionPayloadScript{},
		TransactionPayloadModuleBundle{}, // TODO: ModuleBundle will be removed.
		TransactionPayloadEntryFunction{},
	)

	lcs.RegisterEnum(
		(*TransactionArgument)(nil),

		TransactionArgumentU8{},
		TransactionArgumentU64{},
		TransactionArgumentU128{},
		TransactionArgumentAddress{},
		TransactionArgumentU8Vector{},
		TransactionArgumentBool{},
	)

	lcs.RegisterEnum(
		(*RawTransactionWithData)(nil),

		MultiAgentRawTransaction{},
	)
}

type RawTransaction struct {
	Sender                  AccountAddress     `lcs:"sender"`
	SequenceNumber          uint64             `lcs:"sequence_number"`
	Payload                 TransactionPayload `lcs:"payload"`
	MaxGasAmount            uint64             `lcs:"max_gas_amount"`
	GasUnitPrice            uint64             `lcs:"gas_unit_price"`
	ExpirationTimestampSecs uint64             `lcs:"expiration_timestamp_secs"`
	ChainId                 uint8              `lcs:"chain_id"`
}

type TransactionPayload interface{}

type TransactionPayloadScript struct {
	Code   []byte                `lcs:"code"`
	TyArgs []TypeTag             `lcs:"ty_args"`
	Args   []TransactionArgument `lcs:"args"`
}

type TransactionPayloadEntryFunction struct {
	ModuleName   ModuleId   `lcs:"module_name"`
	FunctionName Identifier `lcs:"function_name"`
	TyArgs       []TypeTag  `lcs:"ty_args"`
	Args         [][]byte   `lcs:"args"`
}

type TransactionPayloadModuleBundle struct{}

type Module struct {
	Code []byte `lcs:"code"`
}

type ModuleId struct {
	Address AccountAddress `lcs:"address"`
	Name    Identifier     `lcs:"name"`
}

func NewModuleIdFromString(moduleId string) (*ModuleId, error) {
	parts := strings.Split(moduleId, "::")
	if len(parts) != 2 {
		return nil, errors.New("Invalid module id.")
	}
	addr, err := NewAccountAddressFromHex(parts[0])
	if err != nil {
		return nil, err
	}
	return &ModuleId{
		*addr,
		Identifier(parts[1]),
	}, nil
}

type TransactionArgument interface{}

type TransactionArgumentU8 struct {
	Value uint8 `lcs:"value"`
}
type TransactionArgumentU64 struct {
	Value uint64 `lcs:"value"`
}
type TransactionArgumentU128 struct {
	Uint128
}
type TransactionArgumentAddress struct {
	Value AccountAddress `lcs:"value"`
}
type TransactionArgumentU8Vector struct {
	Value []uint8 `lcs:"value"`
}
type TransactionArgumentBool struct {
	Value bool `lcs:"value"`
}

type RawTransactionWithData interface{}

type MultiAgentRawTransaction struct {
	RawTransaction           RawTransaction   `lcs:"raw_txn"`
	SecondarySignerAddresses []AccountAddress `lcs:"secondary_signer_addresses"`
}

type SignedTransaction struct {
	Transaction   *RawTransaction          `lcs:"transaction"`
	Authenticator TransactionAuthenticator `lcs:"authenticator"`
}

func GenerateBCSTransaction(from *aptosaccount.Account, txn *RawTransaction) ([]byte, error) {
	builder := NewTransactionBuilderEd25519(func(sm SigningMessage) []byte {
		return from.Sign(sm, "")
	}, from.PublicKey)
	return builder.Sign(txn)
}

func GenerateBCSSimulation(from *aptosaccount.Account, txn *RawTransaction) ([]byte, error) {
	builder := NewTransactionBuilderEd25519(func(sm SigningMessage) []byte {
		zero := [32]byte{}
		privateKey := ed25519.NewKeyFromSeed(zero[:])
		return ed25519.Sign(privateKey, sm)
	}, from.PublicKey)
	return builder.Sign(txn)
}
