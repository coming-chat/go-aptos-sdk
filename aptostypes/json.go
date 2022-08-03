package aptostypes

import (
	"encoding/json"
	"math/big"
	"reflect"
	"strconv"
)

var (
	bigT    = reflect.TypeOf((*Big)(nil))
	uint64T = reflect.TypeOf(Uint64(0))
)

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

// Uint64 marshals/unmarshals as a JSON string
type Uint64 uint64

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint64) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(uint64T)
	}
	res, err := strconv.ParseUint(string(input[1:len(input)-1]), 10, 64)
	if err != nil {
		return toUnmarshalTypeError(err, uint64T)
	}
	*b = Uint64(res)
	return nil
}

func (b Uint64) MarshalJSON() ([]byte, error) {
	res := []byte{'"'}
	res = strconv.AppendUint(res, uint64(b), 10)
	res = append(res, '"')
	return res, nil
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

func errNonString(typ reflect.Type) error {
	return &json.UnmarshalTypeError{Value: "non-string", Type: typ}
}

func toUnmarshalTypeError(err error, typ reflect.Type) error {
	return &json.UnmarshalTypeError{Value: err.Error(), Type: typ}
}
