package nft

import "testing"

func TestFetchGraphQL(t *testing.T) {
	// owner := "0x6ed6f83f1891e02c00c58bf8172e3311c982b1c4fbb1be2d85a55562d4085fb1"
	owner := "0xf5bb1482c28e3c600edf4cac9a10511b3d9a8e162d5b64d9741e2a8cb086bb50"
	tokens, err := FetchGraphqlTokensOfOwner(owner, "")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(tokens)
	}
}
