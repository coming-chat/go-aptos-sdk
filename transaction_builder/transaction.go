package transactionbuilder

import (
	"errors"
	"math/big"
	"strings"

	"github.com/the729/lcs"
)

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
	Value [16]byte `lcs:"value"`
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

func (u *TransactionArgumentU128) BigValue() *big.Int {
	bytes := [16]byte{}
	for i := 0; i < 16; i++ {
		bytes[i] = u.Value[16-i-1]
	}
	return big.NewInt(0).SetBytes(bytes[:])
}
func (u *TransactionArgumentU128) SetBigValue(value *big.Int) error {
	if value.Sign() == -1 {
		return errors.New("Invalid U128: negative number.")
	}
	bytes := value.Bytes()
	l := len(bytes)
	if l > 16 {
		return errors.New("Invalid U128: too large number.")
	}
	for i := 0; i < 16; i++ {
		if i >= l {
			u.Value[i] = 0
		} else {
			u.Value[i] = bytes[l-i-1]
		}
	}
	return nil
}

// type SignedTransaction struct {
// 	RawTransaction
// }

func registerTransaction() {
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
