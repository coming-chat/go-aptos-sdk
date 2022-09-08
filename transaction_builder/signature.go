package transactionbuilder

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"github.com/coming-chat/lcs"
	"golang.org/x/crypto/sha3"
)

const (
	ED25519_PUBLICKEY_LENGTH = 32
	ED25519_SIGNATURE_LENGTH = 64

	MAX_SIGNATURES_SUPPORTED = 32

	MULTI_ED25519_SIGNATURE_BITMAP_LENGTH = 4

	MULTI_ED25519_SCHEME = 0x1
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

func NewMultiEd25519PublicKey(publicKeys [][]byte, threshold uint8) (*MultiEd25519PublicKey, error) {
	if threshold > MAX_SIGNATURES_SUPPORTED {
		return nil, fmt.Errorf(`"threshold" cannot be larger than %v`, MAX_SIGNATURES_SUPPORTED)
	}
	if int(threshold) > len(publicKeys) {
		return nil, errors.New(`"threshold" cannot be larger than public key count.`)
	}
	pubkeys := []Ed25519PublicKey{}
	for _, bytes := range publicKeys {
		pubkey, err := NewEd25519PublicKey(bytes)
		if err != nil {
			return nil, err
		}
		pubkeys = append(pubkeys, *pubkey)
	}
	return &MultiEd25519PublicKey{
		PublicKeys: pubkeys,
		Threshold:  threshold,
	}, nil
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
	bytes := append(mp.ToBytes(), MULTI_ED25519_SCHEME)
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

func NewMultiEd25519Signature(signatures [][]byte, bits []uint8) (*MultiEd25519Signature, error) {
	if len(signatures) != len(bits) {
		return nil, errors.New("The number of signatures and bits are not the same.")
	}
	signs := []Ed25519Signature{}
	sort.Slice(signatures, func(i, j int) bool {
		return bits[i] < bits[j]
	})
	for _, bytes := range signatures {
		signature, err := NewEd25519Signature(bytes)
		if err != nil {
			return nil, err
		}
		signs = append(signs, *signature)
	}
	bitmap, err := CreateBitmap(bits)
	if err != nil {
		return nil, err
	}
	return &MultiEd25519Signature{
		Signatures: signs,
		Bitmap:     bitmap,
	}, nil
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
 * The result bitmap should be 0b10100000000000000000000000000001
 *
 * @returns bitmap that is 32bit long
 */
func CreateBitmap(bits []uint8) ([]byte, error) {
	const firstBitInByte byte = 0b10000000
	bitmap := [4]byte{}
	dupCheckSet := make(map[uint8]bool)

	for _, bit := range bits {
		if bit >= MAX_SIGNATURES_SUPPORTED {
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
