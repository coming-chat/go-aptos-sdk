package transactionbuilder

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/the729/lcs"
)

func TestTransactionArgumentSerilize(t *testing.T) {
	u128 := TransactionArgumentU128{}
	u128.SetBigValue(big.NewInt(1311768467750121216))

	tests := []struct {
		name    string
		val     TransactionArgument
		want    []byte
		wantErr bool
	}{
		{
			name: "U64",
			val:  TransactionArgumentU64{1311768467750121216},
			want: []byte{0x00, 0xEF, 0xCD, 0xAB, 0x78, 0x56, 0x34, 0x12},
		},
		{
			name: "U128",
			val:  u128,
			want: []byte{0x00, 0xEF, 0xCD, 0xAB, 0x78, 0x56, 0x34, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lcs.Marshal(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionArgumentSerilize error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && bytes.Compare(got, tt.want) != 0 {
				t.Errorf("TransactionArgumentSerilize not match got %x, want %x", got, tt.want)
			}
		})
	}
}

func TestTransactionArgumentU128_SetBigValue(t *testing.T) {
	// max U128 ≈ 3.4E38 < 1E39
	tests := []struct {
		name    string
		val     *big.Int
		wantErr bool
	}{
		{
			name: "small number",
			val:  big.NewInt(100),
		},
		{
			name: "normal big number",
			val:  big.NewInt(0).Exp(big.NewInt(10), big.NewInt(30), nil),
		},
		{
			name: "max u128",
			val:  big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil), big.NewInt(1)),
		},
		{
			name:    "negative number",
			val:     big.NewInt(-100),
			wantErr: true,
		},
		{
			name:    "very big number",
			val:     big.NewInt(0).Exp(big.NewInt(10), big.NewInt(40), nil),
			wantErr: true,
		},
	}
	u := TransactionArgumentU128{} // Test multiple reads and writes
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.SetBigValue(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionArgumentU128.SetBigValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				restoreInt := u.BigValue()
				if tt.val.Cmp(restoreInt) != 0 {
					t.Errorf("TransactionArgumentU128.BigValue() restore big int failed, %v -> %v", tt.val, restoreInt)
				}
			}
		})
	}
}
