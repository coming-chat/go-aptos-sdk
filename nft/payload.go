package nft

import (
	"encoding/hex"

	txnBuilder "github.com/coming-chat/go-aptos/transaction_builder"
)

var TOKEN_ABIS = []string{
	// aptos-token/build/AptosToken/abis/token/create_collection_script.abi
	"01186372656174655F636F6C6C656374696F6E5F736372697074000000000000000000000000000000000000000000000000000000000000000305746F6B656E3020637265617465206120656D70747920746F6B656E20636F6C6C656374696F6E207769746820706172616D65746572730005046E616D6507000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67000B6465736372697074696F6E07000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67000375726907000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E6700076D6178696D756D020E6D75746174655F73657474696E670600",
	// aptos-token/build/AptosToken/abis/token/create_token_script.abi
	"01136372656174655F746F6B656E5F736372697074000000000000000000000000000000000000000000000000000000000000000305746F6B656E1D2063726561746520746F6B656E20776974682072617720696E70757473000D0A636F6C6C656374696F6E07000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E6700046E616D6507000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67000B6465736372697074696F6E07000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67000762616C616E636502076D6178696D756D020375726907000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E670015726F79616C74795F70617965655F61646472657373041A726F79616C74795F706F696E74735F64656E6F6D696E61746F720218726F79616C74795F706F696E74735F6E756D657261746F72020E6D75746174655F73657474696E6706000D70726F70657274795F6B6579730607000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67000F70726F70657274795F76616C7565730606010E70726F70657274795F74797065730607000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E6700",
	// aptos-token/build/AptosToken/abis/token/direct_transfer_script.abi
	"01166469726563745f7472616e736665725f736372697074000000000000000000000000000000000000000000000000000000000000000305746f6b656e0000051063726561746f72735f61646472657373040a636f6c6c656374696f6e07000000000000000000000000000000000000000000000000000000000000000106737472696e6706537472696e6700046e616d6507000000000000000000000000000000000000000000000000000000000000000106737472696e6706537472696e67001070726f70657274795f76657273696f6e0206616d6f756e7402",
	// aptos-token/build/AptosToken/abis/token_transfers/offer_script.abi
	"010C6F666665725F73637269707400000000000000000000000000000000000000000000000000000000000000030F746F6B656E5F7472616E7366657273000006087265636569766572040763726561746F72040A636F6C6C656374696F6E07000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E6700046E616D6507000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67001070726F70657274795F76657273696F6E0206616D6F756E7402",
	// aptos-token/build/AptosToken/abis/token_transfers/claim_script.abi
	"010C636C61696D5F73637269707400000000000000000000000000000000000000000000000000000000000000030F746F6B656E5F7472616E73666572730000050673656E646572040763726561746F72040A636F6C6C656374696F6E07000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E6700046E616D6507000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67001070726F70657274795F76657273696F6E02",
	// aptos-token/build/AptosToken/abis/token_transfers/cancel_offer_script.abi
	"011363616E63656C5F6F666665725F73637269707400000000000000000000000000000000000000000000000000000000000000030F746F6B656E5F7472616E7366657273000005087265636569766572040763726561746F72040A636F6C6C656374696F6E07000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E6700046E616D6507000000000000000000000000000000000000000000000000000000000000000106737472696E6706537472696E67001070726F70657274795F76657273696F6E02",
}

const MAX_U64 = ^uint64(0)

type NFTRoyalty struct {
	PayeeAddress      txnBuilder.AccountAddress
	PointsDenominator uint64
	PointsNumerator   uint64
}

type NFTProperty struct {
	keys   []string
	values []string
	types  []string
}

type NFTPayloadBuilder struct {
	builder *txnBuilder.TransactionBuilderABI
}

func NewNFTPayloadBuilder() (*NFTPayloadBuilder, error) {
	abibytes := [][]byte{}
	for _, hexString := range TOKEN_ABIS {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			return nil, err
		}
		abibytes = append(abibytes, bytes)
	}
	builder, err := txnBuilder.NewTransactionBuilderABI(abibytes)
	if err != nil {
		return nil, err
	}
	return &NFTPayloadBuilder{builder}, nil
}

/**
 * Creates a new NFT collection payload.
 *
 * @param name Collection name
 * @param description Collection description
 * @param uri URL to additional info about collection
 * @param maxAmount Maximum number of `token_data` allowed within this collection
 */
func (n *NFTPayloadBuilder) CreateCollection(name, description, uri string, maxAmount uint64) (txnBuilder.TransactionPayload, error) {
	if maxAmount == 0 {
		maxAmount = MAX_U64
	}
	return n.builder.BuildTransactionPayload(
		"0x3::token::create_collection_script",
		[]string{},
		[]any{
			name, description, uri, maxAmount, []bool{false, false, false},
		},
	)
}

/**
 * Creates a new NFT payload.
 *
 * @param collectionName Name of collection, that token belongs to
 * @param name Token name
 * @param description Token description
 * @param supply Token supply
 * @param uri URL to additional info about token
 * @param max The maxium of tokens can be minted from this token
 * @param royalty.PayeeAddress the address to receive the royalty
 * @param royalty.PointsDenominator the denominator for calculating royalty
 * @param royalty.PointsNumerator the numerator for calculating royalty
 * @param property.Keys the property keys for storing on-chain properties
 * @param property.Values the property values to be stored on-chain
 * @param property.Types the type of property values
 */
func (n *NFTPayloadBuilder) CreateToken(collectionName, name, description, uri string, supply, max uint64, royalty NFTRoyalty, property *NFTProperty) (txnBuilder.TransactionPayload, error) {
	if max == 0 {
		max = MAX_U64
	}
	if property == nil {
		property = &NFTProperty{}
	}
	return n.builder.BuildTransactionPayload(
		"0x3::token::create_token_script",
		[]string{},
		[]any{
			collectionName,
			name,
			description,
			supply,
			max,
			uri,
			royalty.PayeeAddress,
			royalty.PointsDenominator,
			royalty.PointsNumerator,
			[]bool{false, false, false, false, false},
			property.keys,
			property.values,
			property.types,
		},
	)
}

/**
 * Offer token payload.
 *
 * @param receiver  Hex-encoded 32 byte Aptos account address to which tokens will be transfered
 * @param creator Hex-encoded 32 byte Aptos account address to which created tokens
 * @param collectionName Name of collection where token is stored
 * @param name Token name
 * @param amount Amount of tokens which will be transfered
 * @param propertyVersion the version of token PropertyMap with a default value 0.
 */
func (n *NFTPayloadBuilder) OfferToken(receiver, creator txnBuilder.AccountAddress, collectionName, name string, amount uint64, propertyVersion uint64) (txnBuilder.TransactionPayload, error) {
	return n.builder.BuildTransactionPayload(
		"0x3::token_transfers::offer_script",
		[]string{},
		[]any{
			receiver, creator, collectionName, name, propertyVersion, amount,
		},
	)
}

/**
 * Claims a token
 *
 * @param sender Hex-encoded 32 byte Aptos account address which holds a token
 * @param creator Hex-encoded 32 byte Aptos account address which created a token
 * @param collectionName Name of collection where token is stored
 * @param name Token name
 * @param propertyVersion the version of token PropertyMap with a default value 0.
 */
func (n *NFTPayloadBuilder) ClaimToken(sender, creator txnBuilder.AccountAddress, collectionName, name string, propertyVersion uint64) (txnBuilder.TransactionPayload, error) {
	return n.builder.BuildTransactionPayload(
		"0x3::token_transfers::claim_script",
		[]string{},
		[]any{
			sender, creator, collectionName, name, propertyVersion,
		},
	)
}

/**
 * Removes a token from pending claims list
 *
 * @param receiver Hex-encoded 32 byte Aptos account address which had to claim token
 * @param creator Hex-encoded 32 byte Aptos account address which created a token
 * @param collectionName Name of collection where token is strored
 * @param name Token name
 * @param propertyVersion the version of token PropertyMap with a default value 0.
 */
func (n *NFTPayloadBuilder) CancelTokenOffer(receiver, creator txnBuilder.AccountAddress, collectionName, name string, propertyVersion uint64) (txnBuilder.TransactionPayload, error) {
	return n.builder.BuildTransactionPayload(
		"0x3::token_transfers::cancel_offer_script",
		[]string{},
		[]any{
			receiver, creator, collectionName, name, propertyVersion,
		},
	)
}

// TODO: Directly transfer the specified amount of tokens from account to receiver
// It's using a single multi signature transaction.
// func (n *NFTPayloadBuilder) DirectTransferToken(sender AptosAccount, receiver: AptosAccount, creator: MaybeHexString, collectionName: string, name: string, amount: number, propertyVersion?: number)
