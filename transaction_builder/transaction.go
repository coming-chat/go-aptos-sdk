package transactionbuilder

import (
	"crypto/ed25519"
	"errors"
	"math/big"
	"strings"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/lcs"
)

func init() {
	lcs.RegisterEnum(
		(*TransactionPayload)(nil),

		TransactionPayloadScript{},
		TransactionPayloadModuleBundle{},
		TransactionPayloadEntryFunction{},
	)

	lcs.RegisterEnum(
		(*TransactionArgument)(nil),

		TransactionArgumentU8{},
		TransactionArgumentU64{},
		TransactionArgumentU128{},
		TransactionArgumentAddress{},
		TransactionArgumentAddress{},
		TransactionArgumentU8Vector{},
		TransactionArgumentBool{},
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

type TransactionPayloadModuleBundle struct {
	Codes []Module `lcs:"codes"`
}

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
	Value *big.Int
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

func (u TransactionArgumentU128) MarshalLCS(e *lcs.Encoder) error {
	if u.Value == nil || u.Value.Sign() == -1 {
		return errors.New("Invalid U128: invalid number.")
	}
	bytes := u.Value.Bytes()
	l := len(bytes)
	if l > 16 {
		return errors.New("Invalid U128: too large number.")
	}
	result := [16]byte{}
	for i := 0; i < 16; i++ {
		if i >= l {
			result[i] = 0
		} else {
			result[i] = bytes[l-i-1]
		}
	}
	return e.EncodeFixedBytes(result[:])
}

func (u *TransactionArgumentU128) UnmarshalLCS(d *lcs.Decoder) error {
	bytes, err := d.DecodeFixedBytes(16)
	if err != nil {
		return err
	}
	// reverse bytes
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	u.Value = big.NewInt(0).SetBytes(bytes)
	return nil
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
