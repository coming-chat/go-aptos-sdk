package transactionbuilder

import (
	"testing"
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
