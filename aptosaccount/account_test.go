package aptosaccount

import (
	"testing"
)

const mnemonic = "crack coil okay hotel glue embark all employ east impact stomach cigar"

func TestAccountSign(t *testing.T) {
	account, err := NewAccountWithMnemonic(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	// 0x1d712fcce859405d768bc636f12d0f8ac5ad88b39178214b22685a9cff310fb6
	// 0x55c15111310a9c107745b1cf80d8d9031f0582a1d21a5eeefa0f6e35c4e2ad74
	// 0xe1c1deec04ed6d7f92f867875c5c9733b64e376ca5a7f5da5b6bdaf3dd28eb9c
	t.Logf("pri = %x", account.PrivateKey[:32])
	t.Logf("pub = %x", account.PublicKey)
	t.Logf("add = %x", account.AuthKey)

	data := []byte{0x1}
	salt := "APTOS::RawTransaction"
	signature := account.Sign(data, salt)

	t.Logf("%x", signature)
}
