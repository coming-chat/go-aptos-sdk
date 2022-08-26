package transactionbuilder

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/the729/lcs"
)

func TestNewAccountAddressFromHex(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{
			addr: "0x1aa",
		},
		{
			addr: "0Xaac4",
		},
		{
			addr: "023",
		},
		{
			addr:    "0213w",
			wantErr: true,
		},
		{
			addr:    "12345678901234567890123456789012345678901234567890123456789012345",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountAddressFromHex(tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountAddressFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t.Logf("result = %x", got)
			}
		})
	}
}

func TestBCSSerializeBasicValue(t *testing.T) {
	checker := func(bcs []byte, x interface{}) {
		b, _ := lcs.Marshal(x)
		if bytes.Compare(b, bcs) != 0 {
			t.Errorf("BCSSerializeBasicValue failed %v", reflect.TypeOf(x).Name())
		}
	}
	{
		x := uint8(100)
		checker(BCSSerializeBasicValue(x), x)
	}
	{
		x := uint32(200)
		checker(BCSSerializeBasicValue(x), x)
	}
	{
		x := int64(300)
		checker(BCSSerializeBasicValue(x), x)
	}
	{
		x := "400abc"
		checker(BCSSerializeBasicValue(x), x)
	}
}
