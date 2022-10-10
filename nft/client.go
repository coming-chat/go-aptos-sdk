package nft

import (
	"encoding/json"
	"strconv"

	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	txnBuilder "github.com/coming-chat/go-aptos/transaction_builder"
)

type CollectionData struct {
	// Describes the collection
	Description string `json:"description"`
	// Unique name within this creators account for this collection
	Name string `json:"name"`
	// URL for additional information/media
	Uri string `json:"uri"`
	// Total number of distinct Tokens tracked by the collection
	Count uint64 `json:"count"`
	// Optional maximum number of tokens allowed within this collections
	Maximum uint64 `json:"maximum"`

	Supply uint64 `json:"supply"`
}

type TokenData struct {
	// Unique name within this creators account for this Token's collection
	Collection string `json:"collection"`
	// Describes this Token
	Description string `json:"description"`
	// The name of this Token
	Name string `json:"name"`
	// Optional maximum number of this type of Token.
	Maximum uint64 `json:"maximum"`
	// Total number of this type of Token
	Supply uint64 `json:"supply"`
	/// URL for additional information / media
	Uri string `json:"uri"`
}

type TokenDataId struct {
	/** Token creator address */
	Creator string `json:"creator"`
	/** Unique name within this creator's account for this Token's collection */
	Collection string `json:"collection"`
	/** Name of Token */
	Name string `json:"name"`
}

func (id *TokenDataId) identifier() string {
	return id.Creator + id.Collection + id.Name
}

type TokenId struct {
	TokenDataId TokenDataId `json:"token_data_id"`
	/** version number of the property map */
	PropertyVersion string `json:"property_version"`
}

type Token struct {
	Id TokenId `json:"id"`
	/** server will return string for u64 */
	Amount string `json:"amount"`
}

type TokenClient struct {
	*aptosclient.RestClient
}

func NewTokenClient(client *aptosclient.RestClient) *TokenClient {
	return &TokenClient{client}
}

/**
 * Queries collection data
 * @param creator Hex-encoded 32 byte Aptos account address which created a collection
 * @param collectionName Collection name
 */
func (c *TokenClient) GetCollectionData(creator txnBuilder.AccountAddress, collectionName string) (*CollectionData, error) {
	collections, err := c.GetAccountResource(creator.ToShortString(), "0x3::token::Collections", 0)
	if err != nil {
		return nil, err
	}

	handle := ""
	if data, ok := collections.Data["collection_data"].(map[string]interface{}); ok {
		handle, _ = data["handle"].(string)
	}
	body := aptosclient.TableItemRequest{
		KeyType:   "0x1::string::String",
		ValueType: "0x3::token::CollectionData",
		Key:       collectionName,
	}

	out := struct {
		CollectionData
		CountString  string `json:"count"`
		MaxString    string `json:"maximum"`
		SupplyString string `json:"supply"`
	}{}
	err = c.GetTableItem(&out, handle, body, "")
	if err != nil {
		return nil, err
	}
	out.Count, _ = strconv.ParseUint(out.CountString, 10, 64)
	out.Maximum, _ = strconv.ParseUint(out.MaxString, 10, 64)
	out.Supply, _ = strconv.ParseUint(out.SupplyString, 10, 64)

	return &out.CollectionData, nil
}

/**
 * Queries token data from collection
 *
 * @param creator Hex-encoded 32 byte Aptos account address which created a token
 * @param collectionName Name of collection, which holds a token
 * @param tokenName Token name
 */
func (c *TokenClient) GetTokenData(creator txnBuilder.AccountAddress, collectionName, tokenName string) (*TokenData, error) {
	collections, err := c.GetAccountResource(creator.ToShortString(), "0x3::token::Collections", 0)
	if err != nil {
		return nil, err
	}

	handle := ""
	if data, ok := collections.Data["token_data"].(map[string]interface{}); ok {
		handle, _ = data["handle"].(string)
	}
	tokenDataId := TokenDataId{
		Creator:    creator.ToShortString(),
		Collection: collectionName,
		Name:       tokenName,
	}
	body := aptosclient.TableItemRequest{
		KeyType:   "0x3::token::TokenDataId",
		ValueType: "0x3::token::TokenData",
		Key:       tokenDataId,
	}

	out := struct {
		TokenData
		MaxString    string `json:"maximum"`
		SupplyString string `json:"supply"`
	}{}
	err = c.GetTableItem(&out, handle, body, "")
	if err != nil {
		return nil, err
	}
	out.Maximum, _ = strconv.ParseUint(out.MaxString, 10, 64)
	out.Supply, _ = strconv.ParseUint(out.SupplyString, 10, 64)
	out.Collection = collectionName

	return &out.TokenData, nil
}

/**
 * Queries token balance for a token account
 * @param account Hex-encoded 32 byte Aptos account address which created a token
 * @param tokenId token id
 */
func (c *TokenClient) GetTokenForAccount(account txnBuilder.AccountAddress, tokenId TokenId) (*Token, error) {
	if tokenId.PropertyVersion == "" {
		tokenId.PropertyVersion = "0"
	}
	tokenStore, err := c.GetAccountResource(account.ToShortString(), "0x3::token::TokenStore", 0)
	if err != nil {
		return nil, err
	}

	handle := ""
	if data, ok := tokenStore.Data["tokens"].(map[string]interface{}); ok {
		handle, _ = data["handle"].(string)
	}
	body := aptosclient.TableItemRequest{
		KeyType:   "0x3::token::TokenId",
		ValueType: "0x3::token::Token",
		Key:       tokenId,
	}

	var out Token
	err = c.GetTableItem(&out, handle, body, "")
	if err != nil {
		if restErr, ok := err.(*aptostypes.RestError); ok && restErr.Code == 404 {
			return &Token{
				Id:     tokenId,
				Amount: "0",
			}, nil
		}
		return nil, err
	}

	return &out, nil
}

type NFTInfo struct {
	TokenData        *TokenData
	TokenId          *TokenDataId
	RelatedHash      string
	RelatedTimestamp uint64
}

func (c *TokenClient) GetAllTokenForAccount(account txnBuilder.AccountAddress) ([]*NFTInfo, error) {
	// 我们需要遍历该用户所有的交易，从中筛选出获得 NFT 的交易，再根据其中的 NFT 信息去查询详细的数据

	owner := account.ToShortString()

	nfts := []*NFTInfo{}
	existsNfts := make(map[string]bool, 0)
	parseNftsFromTransactions := func(txns []aptostypes.Transaction) error {
		for _, txn := range txns {
			if !txn.Success {
				continue
			}
			for _, event := range txn.Events {
				if event.Type != "0x3::token::DepositEvent" || event.Guid.AccountAddress != owner {
					continue
				}
				bytes, err := json.Marshal(event.Data)
				if err != nil {
					continue
				}
				token := Token{}
				err = json.Unmarshal(bytes, &token)
				if err != nil {
					continue
				}
				tokenInfo := token.Id.TokenDataId
				if existsNfts[tokenInfo.identifier()] {
					continue
				}
				creator, err := txnBuilder.NewAccountAddressFromHex(tokenInfo.Creator)
				if err != nil {
					continue
				}

				tokenData, err := c.GetTokenData(*creator, tokenInfo.Collection, tokenInfo.Name)
				if err != nil {
					return err
				}

				nft := NFTInfo{
					TokenData:        tokenData,
					TokenId:          &tokenInfo,
					RelatedHash:      txn.Hash,
					RelatedTimestamp: txn.Timestamp,
				}
				nfts = append(nfts, &nft)
				existsNfts[tokenInfo.identifier()] = true
			}
		}
		return nil
	}

	const limit = 200
	offset := uint64(0)
	for {
		txns, err := c.GetAccountTransactions(owner, offset, limit)
		if err != nil {
			return nil, err
		}

		err = parseNftsFromTransactions(txns)
		if err != nil {
			return nil, err
		}

		if len(txns) < limit {
			break
		} else {
			offset = offset + limit
		}
	}

	return nfts, nil
}
