package transactionbuilder

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/coming-chat/lcs"
	"github.com/stretchr/testify/assert"
)

func TestTypeTagParser_ParseTypeTag(t *testing.T) {
	// test TypeTag bool, u8, u64, u128, address, vector
	TestTypeTagParser_ParseTypeTag_Basic(t)
	// test TypeTag struct
	TestTypeTagParser_ParseTypeTag_Struct(t)
	// test error case
	TestTypeTagParser_ParseTypeTag_ShouldError(t)
}

func TestTypeTagParser_ParseTypeTag_Basic(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    TypeTag
		wantErr bool
	}{
		{
			name: "parses bool",
			tag:  "bool",
			want: TypeTagBool{},
		},
		{
			name: "parses u8",
			tag:  "u8",
			want: TypeTagU8{},
		},
		{
			name: "parses u64",
			tag:  "u64",
			want: TypeTagU64{},
		},
		{
			name: "parses u128",
			tag:  "u128",
			want: TypeTagU128{},
		},
		{
			name: "parses address",
			tag:  "address",
			want: TypeTagAddress{},
		},
		{
			name: "parses address",
			tag:  "address",
			want: TypeTagAddress{},
		},
		{
			name: "parses vector address",
			tag:  "vector<address>",
			want: TypeTagVector{Value: TypeTagAddress{}},
		},
		{
			name: "parses vector u64",
			tag:  "vector<   u64  >",
			want: TypeTagVector{Value: TypeTagU64{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewTypeTagParser(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeTagParser.ParseTypeTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := p.ParseTypeTag()
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeTagParser.ParseTypeTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypeTagParser.ParseTypeTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeTagParser_ParseTypeTag_Struct(t *testing.T) {
	typeTagFrom := func(tagString string) TypeTagStruct {
		parser, err := NewTypeTagParser(tagString)
		assert.Nil(t, err, err)
		tag, err := parser.ParseTypeTag()
		assert.Nil(t, err, err)
		assert.IsType(t, tag, TypeTagStruct{})
		return tag.(TypeTagStruct)
	}

	assertStruct := func(tag TypeTag, shortAddress, moduleName, structName string) {
		t.Run("parse_struct", func(t *testing.T) {
			s := tag.(TypeTagStruct)
			assert.IsType(t, tag, TypeTagStruct{})
			assert.Equal(t, s.Address.ToShortString(), shortAddress)
			assert.Equal(t, s.ModuleName, Identifier(moduleName))
			assert.Equal(t, s.Name, Identifier(structName))
		})
	}

	testCoin := typeTagFrom("0x1::test_coin::Coin")
	assertStruct(testCoin, "0x1", "test_coin", "Coin")

	aptosCoin := typeTagFrom(
		"0x1::coin::CoinStore < 0x1::test_coin::AptosCoin1 ,  0x1::test_coin::AptosCoin2 > ")
	assertStruct(aptosCoin, "0x1", "coin", "CoinStore")

	aptosCoinTrailingComma := typeTagFrom(
		"0x1::coin::CoinStore < 0x1::test_coin::AptosCoin1 ,  0x1::test_coin::AptosCoin2, > ")
	assertStruct(aptosCoinTrailingComma, "0x1", "coin", "CoinStore")

	structTypeTags := aptosCoin.TypeArgs
	assert.Equal(t, len(structTypeTags), 2)
	assertStruct(structTypeTags[0], "0x1", "test_coin", "AptosCoin1")
	assertStruct(structTypeTags[1], "0x1", "test_coin", "AptosCoin2")

	coinComplex := typeTagFrom(
		"0x1::coin::CoinStore < 0x2::coin::LPCoin < 0x1::test_coin::AptosCoin1 <u8>, vector<0x1::test_coin::AptosCoin2 > > >")
	assertStruct(coinComplex, "0x1", "coin", "CoinStore")
	assertStruct(coinComplex.TypeArgs[0], "0x2", "coin", "LPCoin")
}

func TestTypeTagParser_ParseTypeTag_ShouldError(t *testing.T) {
	typeTagShouldError := func(tagString string) {
		t.Run("should error", func(t *testing.T) {
			parser, err := NewTypeTagParser(tagString)
			if err != nil {
				return
			}
			_, err = parser.ParseTypeTag()
			if err != nil {
				return
			}
			t.Fatal("The invalid tag should return an error")
		})
	}

	typeTagShouldError("0x1::test_coin")
	typeTagShouldError("0x1::test_coin::CoinStore<0x1::test_coin::AptosCoin")
	typeTagShouldError("0x1::test_coin::CoinStore<0x1::test_coin>")
	typeTagShouldError("0x1:test_coin::AptosCoin")
	typeTagShouldError("0x!::test_coin::AptosCoin")
	typeTagShouldError("0x1::test_coin::AptosCoin<")
	typeTagShouldError("0x1::test_coin::CoinStore<0x1::test_coin::AptosCoin,")
	typeTagShouldError("")
	typeTagShouldError("0x1::<::CoinStore<0x1::test_coin::AptosCoin,")
	typeTagShouldError("0x1::test_coin::><0x1::test_coin::AptosCoin,")
	typeTagShouldError("u32")
}

func Test_serializeArg(t *testing.T) {
	type args struct {
		argVal  any
		argType TypeTag
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "serialize bool",
			args: args{true, TypeTagBool{}},
			want: []byte{0x1},
		},
		{
			name:    "error bool",
			args:    args{122, TypeTagBool{}},
			wantErr: true,
		},
		{
			name: "serialize u8",
			args: args{uint8(255), TypeTagU8{}},
			want: []byte{0xff},
		},
		{
			name:    "error u8",
			args:    args{10000, TypeTagU8{}},
			wantErr: true,
		},
		{
			name: "serialize u64",
			args: args{uint64(18446744073709551615), TypeTagU64{}},
			want: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name:    "error u64",
			args:    args{"1000000", TypeTagU64{}},
			wantErr: true,
		},
		{
			name: "serialize u128",
			args: args{Uint128{
				big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil), big.NewInt(1))}, TypeTagU128{}},
			want: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name: "serialize u128 big.int",
			args: args{
				big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(64), nil), big.NewInt(1)), TypeTagU128{}},
			want: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "error u128",
			args:    args{1222, TypeTagU128{}},
			wantErr: true,
		},
		{
			name: "serialize account address string",
			args: args{"0x1", TypeTagAddress{}},
			want: AccountAddressFromHex("0x1")[:],
		},
		{
			name: "serialize account address object",
			args: args{*AccountAddressFromHex("0x222"), TypeTagAddress{}},
			want: AccountAddressFromHex("0x222")[:],
		},
		{
			name:    "error account address",
			args:    args{123456, TypeTagAddress{}},
			wantErr: true,
		},

		{
			name: "serialize vector u8 array",
			args: args{[]uint8{255}, TypeTagVector{TypeTagU8{}}},
			want: []byte{1, 0xff},
		},
		{
			name: "serialize vector string",
			args: args{"abc", TypeTagVector{TypeTagU8{}}},
			want: []byte{3, 0x61, 0x62, 0x63},
		},
		{
			name: "serialize vector u8 array",
			args: args{[]uint8{0x61, 0x62, 0x63}, TypeTagVector{TypeTagU8{}}},
			want: []byte{3, 0x61, 0x62, 0x63},
		},
		{
			name:    "error vector",
			args:    args{123456, TypeTagVector{TypeTagU8{}}},
			wantErr: true,
		},

		{
			name: "serialize struct",
			args: args{"abc", TypeTagStruct{*AccountAddressFromHex("0x1"), "string", "String", []TypeTag{}}},
			want: []byte{0x3, 0x61, 0x62, 0x63},
		},
		{
			name:    "error struct",
			args:    args{"abc", TypeTagStruct{*AccountAddressFromHex("0x3"), "token", "Token", []TypeTag{}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg, err := parseValidArg(tt.args.argVal, tt.args.argType)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("parseValidArg() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			got, err := lcs.Marshal(arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValidArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseValidArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argToTransactionArgument(t *testing.T) {
	type args struct {
		argVal  any
		argType TypeTag
	}
	tests := []struct {
		name    string
		args    args
		want    TransactionArgument
		wantErr bool
	}{
		{
			name:    "unrecognized arg types",
			args:    args{12345, "unknown_type"},
			wantErr: true,
		},
		{
			name: "convert bool",
			args: args{true, TypeTagBool{}},
			want: TransactionArgumentBool{true},
		},
		{
			name:    "error bool",
			args:    args{123, TypeTagBool{}},
			wantErr: true,
		},
		{
			name: "convert u8",
			args: args{uint8(123), TypeTagU8{}},
			want: TransactionArgumentU8{123},
		},
		{
			name:    "error u8",
			args:    args{"123", TypeTagU8{}},
			wantErr: true,
		},
		{
			name: "convert u64",
			args: args{uint64(123), TypeTagU64{}},
			want: TransactionArgumentU64{123},
		},
		{
			name:    "error u64",
			args:    args{"u64", TypeTagU64{}},
			wantErr: true,
		},
		{
			name: "convert u128",
			args: args{TransactionArgumentU128{Uint128{big.NewInt(123)}}, TypeTagU128{}},
			want: TransactionArgumentU128{Uint128{big.NewInt(123)}},
		},
		{
			name: "convert u128 big.int",
			args: args{big.NewInt(98765), TypeTagU128{}},
			want: TransactionArgumentU128{Uint128{big.NewInt(98765)}},
		},
		{
			name:    "error u128",
			args:    args{uint64(123), TypeTagU128{}},
			wantErr: true,
		},
		{
			name: "convert account address",
			args: args{*AccountAddressFromHex("0x1"), TypeTagAddress{}},
			want: TransactionArgumentAddress{*AccountAddressFromHex("0x1")},
		},
		{
			name: "convert account address from string",
			args: args{"0x12333", TypeTagAddress{}},
			want: TransactionArgumentAddress{*AccountAddressFromHex("0x12333")},
		},
		{
			name:    "error account address",
			args:    args{123456, TypeTagAddress{}},
			wantErr: true,
		},
		{
			name: "convert vector",
			args: args{[]uint8{1, 2, 3}, TypeTagVector{TypeTagU8{}}},
			want: TransactionArgumentU8Vector{[]uint8{1, 2, 3}},
		},
		{
			name:    "error vector",
			args:    args{"123456", TypeTagVector{TypeTagU8{}}},
			wantErr: true,
		},
		{
			name:    "unsupport struct",
			args:    args{"abc", TypeTagStruct{*AccountAddressFromHex("0x1"), "string", "String", []TypeTag{}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := argToTransactionArgument(tt.args.argVal, tt.args.argType)
			if (err != nil) != tt.wantErr {
				t.Errorf("argToTransactionArgument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("argToTransactionArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func AccountAddressFromHex(hex string) *AccountAddress {
	addr, err := NewAccountAddressFromHex(hex)
	if err != nil {
		panic(err)
	}
	return addr
}
