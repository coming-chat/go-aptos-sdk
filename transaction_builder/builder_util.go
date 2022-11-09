package transactionbuilder

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/coming-chat/lcs"
)

type Token struct {
	Type  string
	Value string
}

func isWhiteSpace(char byte) bool {
	return char == ' ' || char == '\f' || char == '\n' || char == '\r' || char == '\t' || char == '\v'
}

func isValidAlphabetic(char byte) bool {
	return char == '_' || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

// Returns Token and Token byte size
func nextToken(tagStr string, pos int) (*Token, int, error) {
	c := tagStr[pos]
	switch {
	case c == ':':
		if tagStr[pos+1] == ':' {
			return &Token{"COLON", "::"}, 2, nil
		} else {
			break
		}
	case c == '<':
		return &Token{"LT", "<"}, 1, nil
	case c == '>':
		return &Token{"GT", ">"}, 1, nil
	case c == ',':
		return &Token{"COMMA", ","}, 1, nil
	case isWhiteSpace(c):
		var i = pos + 1
		for ; i < len(tagStr) && isWhiteSpace(tagStr[i]); i++ {
		}
		return &Token{"SPACE", tagStr[pos:i]}, i - pos, nil
	case isValidAlphabetic(c):
		var i = pos + 1
		for ; i < len(tagStr) && isValidAlphabetic(tagStr[i]); i++ {
		}
		return &Token{"IDENT", tagStr[pos:i]}, i - pos, nil
	}

	return nil, 0, errors.New("Unrecognized token.")
}

func tokenize(tagStr string) ([]Token, error) {
	pos := 0
	tokens := []Token{}
	for pos < len(tagStr) {
		token, size, err := nextToken(tagStr, pos)
		if err != nil {
			return nil, err
		}
		if token.Type != "SPACE" {
			tokens = append(tokens, *token)
		}
		pos += size
	}
	return tokens, nil
}

type TypeTagParser struct {
	Tokens []Token
}

var errInvalidTypeTag = errors.New("Invalid type tag.")

func NewTypeTagParser(tyArg string) (*TypeTagParser, error) {
	tokens, err := tokenize(tyArg)
	if err != nil {
		return nil, err
	}
	return &TypeTagParser{Tokens: tokens}, nil
}

func (p *TypeTagParser) shift() *Token {
	if len(p.Tokens) == 0 {
		return nil
	}
	t := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	return &t
}

func (p *TypeTagParser) consume(targetToken string) error {
	token := p.shift()
	if token == nil || token.Value != targetToken {
		return errInvalidTypeTag
	}
	return nil
}

func (p *TypeTagParser) parseCommaList(endToken string, allowTraillingComma bool) ([]TypeTag, error) {
	if len(p.Tokens) <= 0 {
		return nil, errInvalidTypeTag
	}
	res := []TypeTag{}

	for p.Tokens[0].Value != endToken {
		tag, err := p.ParseTypeTag()
		if err != nil {
			return nil, err
		}
		res = append(res, tag)

		if len(p.Tokens) > 0 && p.Tokens[0].Value == endToken {
			break
		}
		err = p.consume(",")
		if err != nil {
			return nil, err
		}
		if len(p.Tokens) > 0 && p.Tokens[0].Value == endToken && allowTraillingComma {
			break
		}

		if len(p.Tokens) <= 0 {
			return nil, errInvalidTypeTag
		}
	}

	return res, nil
}

func (p *TypeTagParser) ParseTypeTag() (TypeTag, error) {
	if len(p.Tokens) == 0 {
		return nil, errInvalidTypeTag
	}
	var err error

	token := p.shift()
	switch token.Value {
	case "u8":
		return TypeTagU8{}, nil
	case "u64":
		return TypeTagU64{}, nil
	case "u128":
		return TypeTagU128{}, nil
	case "bool":
		return TypeTagBool{}, nil
	case "address":
		return TypeTagAddress{}, nil
	case "vector":
		err = p.consume("<")
		if err != nil {
			return nil, err
		}
		tag, err := p.ParseTypeTag()
		if err != nil {
			return nil, err
		}
		err = p.consume(">")
		if err != nil {
			return nil, err
		}
		return TypeTagVector{Value: tag}, nil
	}
	if token.Type == "IDENT" && (strings.HasPrefix(token.Value, "0x") || strings.HasPrefix(token.Value, "0X")) {
		address := token.Value

		err = p.consume("::")
		if err != nil {
			return nil, err
		}
		moduleToken := p.shift()
		if moduleToken == nil || moduleToken.Type != "IDENT" {
			return nil, errInvalidTypeTag
		}

		err = p.consume("::")
		if err != nil {
			return nil, err
		}
		nameToken := p.shift()
		if nameToken == nil || nameToken.Type != "IDENT" {
			return nil, errInvalidTypeTag
		}

		tyArgs := []TypeTag{}
		// Check if the struct has ty args
		if len(p.Tokens) > 0 && p.Tokens[0].Value == "<" {
			err = p.consume("<")
			if err != nil {
				return nil, err
			}
			tyArgs, err = p.parseCommaList(">", true)
			if err != nil {
				return nil, err
			}
			err = p.consume(">")
			if err != nil {
				return nil, err
			}
		}

		addr, err := NewAccountAddressFromHex(address)
		if err != nil {
			return nil, err
		}
		return TypeTagStruct{
			Address:    *addr,
			ModuleName: Identifier(moduleToken.Value),
			Name:       Identifier(nameToken.Value),
			TypeArgs:   tyArgs,
		}, nil
	}
	return nil, errInvalidTypeTag
}

func serializeArg(argVal any, argType TypeTag, encoder *lcs.Encoder) error {
	switch argType.(type) {
	case TypeTagBool:
		if v, ok := argVal.(bool); ok {
			return encoder.Encode(v)
		}
	case TypeTagU8:
		if v, ok := argVal.(uint8); ok {
			return encoder.Encode(v)
		}
		if v, ok := argVal.(int); ok && v == int(uint8(v)) {
			return encoder.Encode(uint8(v))
		}
		if v, ok := argVal.(float64); ok && v == float64(uint8(v)) {
			return encoder.Encode(uint8(v))
		}
		if v, ok := argVal.(string); ok {
			u, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return err
			}
			return encoder.Encode(uint8(u))
		}
	case TypeTagU64:
		if v, ok := argVal.(uint64); ok {
			return encoder.Encode(v)
		}
		if v, ok := argVal.(int); ok && v >= 0 {
			return encoder.Encode(uint64(v))
		}
		if v, ok := argVal.(float64); ok && v >= 0 {
			return encoder.Encode(uint64(v))
		}
		if v, ok := argVal.(string); ok {
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			return encoder.Encode(u)
		}
	case TypeTagU128:
		if v, ok := argVal.(Uint128); ok {
			return encoder.Encode(v)
		}
		if v, ok := argVal.(*big.Int); ok {
			return encoder.Encode(Uint128{v})
		}
		if v, ok := argVal.(int); ok && v >= 0 {
			return encoder.Encode(Uint128{big.NewInt(int64(v))})
		}
		if v, ok := argVal.(float64); ok && v >= 0 {
			return encoder.Encode(Uint128{big.NewInt(int64(v))})
		}
		if v, ok := argVal.(string); ok {
			if big, ok := big.NewInt(0).SetString(v, 10); ok {
				return encoder.Encode(Uint128{big})
			}
		}
	case TypeTagAddress:
		if v, ok := argVal.(AccountAddress); ok {
			return encoder.Encode(v)
		}
		if v, ok := argVal.(string); ok {
			addr, err := NewAccountAddressFromHex(v)
			if err != nil {
				return err
			}
			return encoder.Encode(addr)
		}
	case TypeTagVector:
		itemType := argType.(TypeTagVector).Value
		switch itemType.(type) {
		case TypeTagU8:
			if v, ok := argVal.([]byte); ok {
				return encoder.Encode(v)
			}
			if v, ok := argVal.(string); ok {
				return encoder.Encode(v)
			}
		}

		rv := reflect.ValueOf(argVal)
		kindstring := rv.Kind().String()
		print(kindstring)
		if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
			return errors.New("Invalid vector args.")
		}
		length := rv.Len()
		if err := encoder.EncodeUleb128(uint64(length)); err != nil {
			return err
		}
		for i := 0; i < length; i++ {
			if err := serializeArg(rv.Index(i).Interface(), itemType, encoder); err != nil {
				return err
			}
		}
		return nil
	case TypeTagStruct:
		tag := argType.(TypeTagStruct)
		if tag.ShortFunctionName() != "0x1::string::String" {
			return errors.New("The only supported struct arg is of type 0x1::string::String")
		}
		if v, ok := argVal.(string); ok {
			return encoder.Encode(v)
		}
	default:
		return errors.New("Unsupported arg type.")
	}
	return fmt.Errorf("Invalid argument %v.", argVal)
}

func argToTransactionArgument(argVal any, argType TypeTag) (TransactionArgument, error) {
	switch argType.(type) {
	case TypeTagBool:
		if v, ok := argVal.(bool); ok {
			return TransactionArgumentBool{v}, nil
		}
	case TypeTagU8:
		if v, ok := argVal.(uint8); ok {
			return TransactionArgumentU8{v}, nil
		}
	case TypeTagU64:
		if v, ok := argVal.(uint64); ok {
			return TransactionArgumentU64{v}, nil
		}
	case TypeTagU128:
		if v, ok := argVal.(TransactionArgumentU128); ok {
			return v, nil
		}
		if v, ok := argVal.(Uint128); ok {
			return TransactionArgumentU128{v}, nil
		}
		if v, ok := argVal.(*big.Int); ok {
			return TransactionArgumentU128{Uint128{v}}, nil
		}
	case TypeTagAddress:
		if v, ok := argVal.(AccountAddress); ok {
			return TransactionArgumentAddress{v}, nil
		}
		if s, ok := argVal.(string); ok {
			addr, err := NewAccountAddressFromHex(s)
			if err == nil {
				return TransactionArgumentAddress{*addr}, nil
			}
		}
	case TypeTagVector:
		itemValue := argType.(TypeTagVector).Value
		switch itemValue.(type) {
		case TypeTagU8:
			if v, ok := argVal.([]byte); ok {
				return TransactionArgumentU8Vector{v}, nil
			}
		}
	case TypeTagSigner: // unsupport
		return nil, errors.New("Unknown type for TransactionArgument.")
	case TypeTagStruct: // unsupport
		return nil, errors.New("Unknown type for TransactionArgument.")
	default:
		return nil, errors.New("Unknown type for TransactionArgument.")
	}

	return nil, fmt.Errorf("Invalid argument %v.", argVal)
}
