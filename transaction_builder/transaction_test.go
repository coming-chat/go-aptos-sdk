package transactionbuilder

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/coming-chat/lcs"
)

func TestTransactionArgumentSerilize(t *testing.T) {
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
			val:  TransactionArgumentU128{big.NewInt(1311768467750121216)},
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

func TestTransactionArgumentU128(t *testing.T) {
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
			val:  big.NewInt(0).Exp(big.NewInt(10), big.NewInt(30), nil), // 10 ^ 30
		},
		{
			name: "max u128",
			val:  big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil), big.NewInt(1)), // 2 ^ 128 - 1
		},
		{
			name:    "negative number",
			val:     big.NewInt(-100),
			wantErr: true,
		},
		{
			name:    "very big number",
			val:     big.NewInt(0).Exp(big.NewInt(10), big.NewInt(40), nil), // 10 ^ 40
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := TransactionArgumentU128{tt.val}
			bytes, err := lcs.Marshal(u)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionArgumentU128 marshal error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				u2 := TransactionArgumentU128{}
				err = lcs.Unmarshal(bytes, &u2)
				if err != nil {
					t.Errorf("TransactionArgumentU128 unmarshal error = %v", err)
					return
				}
				if u.Value.Cmp(u2.Value) != 0 {
					t.Errorf("TransactionArgumentU128 restore big int failed, %v -> %v", u.Value, u2.Value)
				}
			}
		})
	}
}
