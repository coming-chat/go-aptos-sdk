package transactionbuilder

import (
	"encoding/hex"
	"errors"
	"strings"
)

const AddressLength = 32

type AccountAddress [AddressLength]byte

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
	if len(bytes) > AddressLength {
		return nil, errors.New("Hex string is too long. Address's length is 32 bytes.")
	}

	res := AccountAddress{}
	copy(res[AddressLength-len(bytes):], bytes[:])
	return &res, nil
}
