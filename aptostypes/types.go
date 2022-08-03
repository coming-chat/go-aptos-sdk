package aptostypes

import (
	"math/big"
)

//go:generate go run github.com/fjl/gencodec -type LedgerInfo -field-override ledgerInfoMarshaling -out gen_ledger_json.go

// LedgerInfo represents chain current ledger info
type LedgerInfo struct {
	ChainId         int      `json:"chain_id"`
	LedgerVersion   *big.Int `json:"ledger_version" gencodec:"required"`
	LedgerTimestamp *big.Int `json:"ledger_timestamp" gencodec:"required"`

	Epoch    int    `json:"epoch"`
	NodeRole string `json:"node_role"`
}

type ledgerInfoMarshaling struct {
	LedgerVersion   *Big
	LedgerTimestamp *Big
}
