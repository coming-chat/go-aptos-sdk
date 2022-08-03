package aptostypes

import (
	"encoding/json"
	"math/big"
	"reflect"
)

var bigT = reflect.TypeOf((*Big)(nil))

// Big represents a string number, eg "1659510704301760"
type Big big.Int

// UnmarshalJSON implements json.Unmarshaler.
func (b *Big) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(bigT)
	}
	res, bo := big.NewInt(0).SetString(string(input[1:len(input)-1]), 10)
	if !bo {
		return errNonString(bigT)
	}
	*b = (Big)(*res)
	return nil
}

func (b *Big) MarshalJSON() ([]byte, error) {
	res := []byte{'"'}
	res = append(res, []byte((*big.Int)(b).String())...)
	res = append(res, '"')
	return res, nil
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

func errNonString(typ reflect.Type) error {
	return &json.UnmarshalTypeError{Value: "non-string", Type: typ}
}
