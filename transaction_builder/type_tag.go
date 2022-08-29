package transactionbuilder

import (
	"errors"
	"strings"

	"github.com/coming-chat/lcs"
)

func init() {
	lcs.RegisterEnum(
		(*TypeTag)(nil),

		TypeTagBool{},
		TypeTagU8{},
		TypeTagU64{},
		TypeTagU128{},
		TypeTagAddress{},
		TypeTagSigner{},
		TypeTagVector{},
		TypeTagStruct{},
	)
}

type TypeTag interface{}

type TypeTagBool struct{}
type TypeTagU8 struct{}
type TypeTagU64 struct{}
type TypeTagU128 struct{}
type TypeTagAddress struct{}
type TypeTagSigner struct{}
type TypeTagVector struct {
	Value TypeTag `lcs:"value"`
}
type TypeTagStruct struct {
	Address    AccountAddress `lcs:"address"`
	ModuleName Identifier     `lcs:"module_name"`
	Name       Identifier     `lcs:"name"`
	TypeArgs   []TypeTag      `lcs:"type_args"`
}

func NewTypeTagStructFromString(tag string) (*TypeTagStruct, error) {
	if strings.Contains(tag, "<") {
		return nil, errors.New("Not implemented")
	}

	parts := strings.Split(tag, "::")
	if len(parts) != 3 {
		return nil, errors.New("Invalid struct tag string literal.")
	}
	addr, err := NewAccountAddressFromHex(parts[0])
	if err != nil {
		return nil, err
	}
	return &TypeTagStruct{
		Address:    *addr,
		ModuleName: Identifier(parts[1]),
		Name:       Identifier(parts[2]),
	}, nil
}
