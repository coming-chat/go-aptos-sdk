package nft

import (
	"github.com/coming-chat/go-aptos/graphql"
	"testing"
)

func TestFetchGraphQL(t *testing.T) {
	owner := "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
	createrAddress := "0x01a97439554e302ee7584f55ea0035b6d4f131e09a024aba121139682a7ab3e8"
	tokens, err := FetchGraphqlTokensOfOwner(owner, graphql.GraphUrlTestnet, "")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(tokens)
	}
	t.Log("<===================test filter====================>")
	tokensWithFilter, err := FetchGraphqlTokensOfOwner(owner, graphql.GraphUrlTestnet, createrAddress)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(tokensWithFilter)
	}
}
