package transactionbuilder

import (
	"context"
	"testing"

	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/stretchr/testify/require"
)

func TestTransactionBuilderRemoteABI(t *testing.T) {
	chain, err := aptosclient.Dial(context.Background(), MainnetRestUrl)
	require.Nil(t, err)

	// callMsg := map[string]any{
	function := "0x89576037b3cc0b89645ea393a47787bb348272c76d6941c574b053672b848039::aggregator::three_step_route"
	type_arguments := []string{"0x1::aptos_coin::AptosCoin", "0x1000000fa32d122c18a6a31c009ce5e71674f22d06a581bb0a15575e6addadcc::usda::USDA", "0x84d7aeef42d38a5ffc3ccef853e1b82e4958659d16a7de736a29c55fbbeb0114::staked_aptos_coin::StakedAptosCoin", "0x5e156f1207d0ebfa19a9eeff00d62a282278fb8719f4fab3a586a0a2c0fffbea::coin::T", "u8", "0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::curves::Uncorrelated", "u8"}

	builder, err := NewTransactionBuilderRemoteABIWithFunc("0x89576037b3cc0b89645ea393a47787bb348272c76d6941c574b053672b848039::aggregator", chain)
	require.Nil(t, err)

	arguments1 := []any{9, 0, true, 3, 0, true, 8, 0, false, 10000000, 694840}
	payload1, err := builder.BuildTransactionPayload(function, type_arguments, arguments1)
	require.Nil(t, err)

	arguments2 := []any{"9", "0", true, 3, 0, true, "8", "0", false, "10000000", 694840} // some u8/u64 are set to strings.
	payload2, err := builder.BuildTransactionPayload(function, type_arguments, arguments2)
	require.Nil(t, err)

	require.Equal(t, payload1, payload2)
}

func TestIntToUintParse(t *testing.T) {
	tests := []struct {
		name  string
		value int
		isU8  bool
	}{
		{
			name:  "negative",
			value: -10,
		},
		{
			name:  "greater than max u8",
			value: 0,
			isU8:  true,
		},
		{
			name:  "max u8",
			value: 255,
			isU8:  true,
		},
		{
			name:  "greater than max u8",
			value: 256,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.isU8, tt.value == int(uint8(tt.value)))
		})
	}
}
