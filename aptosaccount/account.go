package aptosaccount

import (
	"crypto/ed25519"
	"errors"
	"github.com/coming-chat/go-aptos/crypto/derivation"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

type Account struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	AuthKey    [32]byte
}

func NewAccount(seed []byte) *Account {
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)
	data := append(publicKey, 0x00)
	authKey := sha3.Sum256(data)
	return &Account{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AuthKey:    authKey,
	}
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := derivation.DeriveForPath("m/44'/637'/0'/0'/0'", seed)
	if err != nil {
		return nil, err
	}
	return NewAccount(key.Key), nil
}

func GetOldVersionPrivateKeyWithMnemonic(mnemonic string) ([]byte, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	return seed[:32], nil
}

func GenerateMultisignerAuthKey(publicKeys [][]byte, threshold int) ([32]byte, error) {
	N := len(publicKeys)
	if threshold > N {
		return [32]byte{}, errors.New("The threshold must be less than the number of public keys")
	}

	data := []byte{}
	for i := 0; i < N; i++ {
		publicKey := publicKeys[i]
		data = append(data, publicKey...)
	}
	data = append(data, byte(threshold))
	data = append(data, 0x01)
	authKey := sha3.Sum256(data)

	return authKey, nil
}

func (a *Account) Sign(data []byte, salt string) []byte {
	return Sign(a.PrivateKey, data, salt)
}

func Sign(privateKey ed25519.PrivateKey, data []byte, salt string) []byte {
	prefixBytes := []byte{}
	if len(salt) > 0 {
		s := sha3.Sum256([]byte(salt))
		prefixBytes = s[:]
	}

	signingMessage := append(prefixBytes, data...)
	return ed25519.Sign(privateKey, signingMessage)
}
