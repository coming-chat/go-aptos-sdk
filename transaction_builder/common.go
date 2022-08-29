package transactionbuilder

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/coming-chat/lcs"
)

const (
	ADDRESS_LENGTH = 32
)

type AccountAddress [ADDRESS_LENGTH]byte

type Identifier string

/**
 * Creates AccountAddress from a hex string.
 * @param addr Hex string can be with a prefix or without a prefix,
 *   e.g. '0x1aa' or '1aa'. Hex string will be left padded with 0s if too short.
 */
func NewAccountAddressFromHex(addr string) (*AccountAddress, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	if len(bytes) > ADDRESS_LENGTH {
		return nil, fmt.Errorf("Hex string is too long. Address's length is %v bytes.", ADDRESS_LENGTH)
	}

	res := AccountAddress{}
	copy(res[ADDRESS_LENGTH-len(bytes):], bytes[:])
	return &res, nil
}

func BCSSerializeBasicValue[T bool | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | string](t T) []byte {
	s, _ := lcs.Marshal(t)
	return s
}
