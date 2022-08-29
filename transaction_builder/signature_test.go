package transactionbuilder

import (
	"reflect"
	"testing"

	"github.com/coming-chat/lcs"
)

func TestMultiEd25519PublicKey(t *testing.T) {
	pubkey1 := [ED25519_PUBLICKEY_LENGTH]byte{0x1, 0x2, 0x3}
	pubkey2 := [ED25519_PUBLICKEY_LENGTH]byte{0xa, 0xb, 0xc}
	mp := MultiEd25519PublicKey{
		PublicKeys: []Ed25519PublicKey{
			{pubkey1[:]},
			{pubkey2[:]},
		},
		Threshold: 2,
	}

	bytes, err := lcs.Marshal(mp)
	if err != nil {
		t.Fatal(err)
	}

	mp2 := MultiEd25519PublicKey{}
	err = lcs.Unmarshal(bytes, &mp2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(mp2.Threshold, mp.PublicKeys)
}

func TestMultiEd25519Signature(t *testing.T) {
	sign1 := [ED25519_SIGNATURE_LENGTH]byte{0x1, 0x2, 0x3}
	sign2 := [ED25519_SIGNATURE_LENGTH]byte{0xa, 0xb, 0xc}
	ms := MultiEd25519Signature{
		Signatures: []Ed25519Signature{
			{sign1[:]},
			{sign2[:]},
		},
		Bitmap: []byte{0x1, 0x2, 0xa, 0xb},
	}

	bytes, err := lcs.Marshal(ms)
	if err != nil {
		t.Fatal(err)
	}

	ms2 := MultiEd25519Signature{}
	err = lcs.Unmarshal(bytes, &ms2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ms2.Bitmap, ms2.Signatures)
}

func TestCreateBitmap(t *testing.T) {
	tests := []struct {
		name    string
		bits    []uint8
		want    []byte
		wantErr bool
	}{
		{
			bits: []uint8{0, 2, 31},
			want: []byte{0b10100000, 0b00000000, 0b00000000, 0b00000001},
		},
		{
			bits:    []uint8{32},
			wantErr: true,
		},
		{
			bits:    []uint8{2, 2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateBitmap(tt.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBitmap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateBitmap() = %v, want %v", got, tt.want)
			}
		})
	}
}
