package transactionbuilder

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/coming-chat/lcs"
	"golang.org/x/crypto/sha3"
)

const (
	ED25519_PUBLICKEY_LENGTH = 32
	ED25519_SIGNATURE_LENGTH = 64

	MAX_SIGNATURE_SUPPORTED = 32

	MULTI_ED25519_SIGNATURE_BITMAP_LENGTH = 4
)

type PublicKey interface{}
type Signature interface{}

// ------ Ed25519 ------

type Ed25519PublicKey struct {
	PublicKey []byte `lcs:"publicKey"`
}

func NewEd25519PublicKey(publicKey []byte) (*Ed25519PublicKey, error) {
	if len(publicKey) != ED25519_PUBLICKEY_LENGTH {
		return nil, fmt.Errorf(`Ed25519PublicKey length should be %d`, ED25519_PUBLICKEY_LENGTH)
	}
	return &Ed25519PublicKey{publicKey}, nil
}

type Ed25519Signature struct {
	Signature []byte `lcs:"signature"`
}

func NewEd25519Signature(signature []byte) (*Ed25519Signature, error) {
	if len(signature) != ED25519_SIGNATURE_LENGTH {
		return nil, fmt.Errorf(`Ed25519Signature length should be %d`, ED25519_SIGNATURE_LENGTH)
	}
	return &Ed25519Signature{signature}, nil
}

// ------ Multi-Ed25519 ------

type MultiEd25519PublicKey struct {
	PublicKeys []Ed25519PublicKey `lcs:"publicKeys"`
	Threshold  uint8              `lcs:"threshold"`
}

func (mp *MultiEd25519PublicKey) ToBytes() []byte {
	bytes := make([]byte, len(mp.PublicKeys)*ED25519_PUBLICKEY_LENGTH+1)
	for idx, pubkey := range mp.PublicKeys {
		copy(bytes[idx*ED25519_PUBLICKEY_LENGTH:], pubkey.PublicKey)
	}
	bytes[len(mp.PublicKeys)*ED25519_PUBLICKEY_LENGTH] = mp.Threshold
	return bytes
}

func (mp MultiEd25519PublicKey) MarshalLCS(e *lcs.Encoder) error {
	return e.EncodeBytes(mp.ToBytes())
}

func (mp *MultiEd25519PublicKey) UnmarshalLCS(d *lcs.Decoder) error {
	bytes, err := d.DecodeBytes()
	if err != nil {
		return err
	}
	mp.Threshold = bytes[len(bytes)-1]
	mp.PublicKeys = []Ed25519PublicKey{}
	for i := 0; i < len(bytes)-1; i += ED25519_PUBLICKEY_LENGTH {
		publicBytes := bytes[i : i+ED25519_PUBLICKEY_LENGTH]
		mp.PublicKeys = append(mp.PublicKeys, Ed25519PublicKey{publicBytes})
	}
	return nil
}

func (mp *MultiEd25519PublicKey) AuthenticationKey() [32]byte {
	bytes := append(mp.ToBytes(), 0x01)
	authKey := sha3.Sum256(bytes)
	return authKey
}

func (mp *MultiEd25519PublicKey) Address() string {
	b := mp.AuthenticationKey()
	return "0x" + hex.EncodeToString(b[:])
}

type MultiEd25519Signature struct {
	Signatures []Ed25519Signature
	Bitmap     []byte
}

func (ms *MultiEd25519Signature) ToBytes() []byte {
	bytes := make([]byte, len(ms.Signatures)*ED25519_SIGNATURE_LENGTH+MULTI_ED25519_SIGNATURE_BITMAP_LENGTH)
	for idx, signature := range ms.Signatures {
		copy(bytes[idx*ED25519_SIGNATURE_LENGTH:], signature.Signature)
	}
	copy(bytes[len(ms.Signatures)*ED25519_SIGNATURE_LENGTH:], ms.Bitmap)
	return bytes
}

func (ms MultiEd25519Signature) MarshalLCS(e *lcs.Encoder) error {
	return e.EncodeBytes(ms.ToBytes())
}

func (ms *MultiEd25519Signature) UnmarshalLCS(d *lcs.Decoder) error {
	bytes, err := d.DecodeBytes()
	if err != nil {
		return err
	}
	ms.Bitmap = bytes[len(bytes)-MULTI_ED25519_SIGNATURE_BITMAP_LENGTH:]
	ms.Signatures = []Ed25519Signature{}
	for i := 0; i < len(bytes)-MULTI_ED25519_SIGNATURE_BITMAP_LENGTH; i += ED25519_SIGNATURE_LENGTH {
		signatureBytes := bytes[i : i+ED25519_SIGNATURE_LENGTH]
		ms.Signatures = append(ms.Signatures, Ed25519Signature{signatureBytes})
	}
	return nil
}

/**
 * Helper method to create a bitmap out of the specified bit positions
 * @param bits The bitmap positions that should be set. A position starts at index 0.
 * Valid position should range between 0 and 31.
 * @example
 * Here's an example of valid `bits`
 * ```
 * [0, 2, 31]
 * ```
 * `[0, 2, 31]` means the 1st, 3rd and 32nd bits should be set in the bitmap.
 * The result bitmap should be 0b1010000000000000000000000000001
 *
 * @returns bitmap that is 32bit long
 */
func CreateBitmap(bits []uint8) ([]byte, error) {
	const firstBitInByte byte = 0b10000000
	bitmap := [4]byte{}
	dupCheckSet := make(map[uint8]bool)

	for _, bit := range bits {
		if bit >= MAX_SIGNATURE_SUPPORTED {
			return nil, fmt.Errorf("Invalid bit value %v.", bit)
		}
		if dupCheckSet[bit] {
			return nil, errors.New("Duplicated bits detected.")
		}
		dupCheckSet[bit] = true

		byteOffset := bit / 8
		_byte := bitmap[byteOffset]
		_byte |= firstBitInByte >> (bit % 8)
		bitmap[byteOffset] = _byte
	}
	return bitmap[:], nil
}
