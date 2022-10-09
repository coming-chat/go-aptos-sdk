package aptosaccount

import (
	"crypto/ed25519"
	"reflect"
	"testing"
)

const mnemonic = "crack coil okay hotel glue embark all employ east impact stomach cigar"

var seed = [32]byte{164, 52, 187, 8, 138, 232, 166, 157, 88, 132, 167, 31, 232, 86, 153, 185, 160, 88, 207, 158, 43, 104, 143, 80, 48, 155, 175, 178, 241, 78, 196, 78}

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

func TestAccount_Sign_Verify(t *testing.T) {
	type fields struct {
		PrivateKey ed25519.PrivateKey
		PublicKey  ed25519.PublicKey
		AuthKey    [32]byte
	}
	type args struct {
		data []byte
		salt string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test case 1",
			fields: fields{
				PublicKey:  NewAccount(seed[:]).PublicKey,
				PrivateKey: NewAccount(seed[:]).PrivateKey,
				AuthKey:    NewAccount(seed[:]).AuthKey,
			},
			args: args{
				data: []byte{0x01},
				salt: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				PrivateKey: tt.fields.PrivateKey,
				PublicKey:  tt.fields.PublicKey,
				AuthKey:    tt.fields.AuthKey,
			}
			if got := a.Sign(tt.args.data, tt.args.salt); !Verify(a.PublicKey, tt.args.data, got) {
				t.Errorf("Sign() = %v, verify %v", got, false)
			}
		})
	}
}

func TestNewAccountWithMnemonic(t *testing.T) {
	type args struct {
		mnemonic string
	}
	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				mnemonic: mnemonic,
			},
			want:    NewAccount(seed[:]),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.args.mnemonic)
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got, tt.want)
			}
		})
	}
}
