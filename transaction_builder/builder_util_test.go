package transactionbuilder

import (
	"reflect"
	"testing"

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
	t.Log("ParseTypeTag_Basic Pass")
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
		s := tag.(TypeTagStruct)
		assert.IsType(t, tag, TypeTagStruct{})
		assert.Equal(t, s.Address.ToShortString(), shortAddress)
		assert.Equal(t, s.ModuleName, Identifier(moduleName))
		assert.Equal(t, s.Name, Identifier(structName))
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

	t.Log("ParseTypeTag_Struct Pass")
}

func TestTypeTagParser_ParseTypeTag_ShouldError(t *testing.T) {
	typeTagShouldError := func(tagString string) {
		parser, err := NewTypeTagParser(tagString)
		if err != nil {
			return
		}
		_, err = parser.ParseTypeTag()
		if err != nil {
			return
		}
		t.Fatal("The invalid tag should return an error")
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

	t.Log("ParseTypeTag_ShouldError Pass")
}
