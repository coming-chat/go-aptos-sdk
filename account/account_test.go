package aptosaccount

import (
	"os"
	"testing"
)

func TestAccountSign(t *testing.T) {
	mnemonic := os.Getenv("WalletSdkTestM1")
	account, err := NewAccountWithMnemonic(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte{0x1}
	salt := "APTOS::RawTransaction"
	signature := account.Sign(data, salt)

	t.Logf("%x", signature)
}
