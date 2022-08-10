package aptostypes

import "unsafe"

const (
	TypePendingTransaction            = "pending_transaction"
	TypeGenesisTransaction            = "genesis_transaction"
	TypeUserTransaction               = "user_transaction"
	TypeBlockMetadataTransaction      = "block_metadata_transaction"
	TypeStateCheckoutpointTransaction = "state_checkpoint_transaction"
)

// use unsafe.Pointer transfer Transaction to UserTransaction,BlockMetadataTransaction
// the fields order and types must be consistent between Transaction and or XXTransaction

type UserTransaction struct {
	Type                    string     `json:"type"`
	Hash                    string     `json:"hash"`
	Sender                  string     `json:"sender"`
	SequenceNumber          uint64     `json:"sequence_number"`
	MaxGasAmount            uint64     `json:"max_gas_amount"`
	GasUnitPrice            uint64     `json:"gas_unit_price"`
	GasCurrencyCode         string     `json:"gas_currency_code"`
	ExpirationTimestampSecs uint64     `json:"expiration_timestamp_secs"`
	Payload                 *Payload   `json:"payload"`
	Signature               *Signature `json:"signature"`
	Events                  []Event    `json:"events"`
	Version                 uint64     `json:"version"`
	StateRootHash           string     `json:"state_root_hash"`
	EventRootHash           string     `json:"event_root_hash"`
	GasUsed                 uint64     `json:"gas_used"`
	Success                 bool       `json:"success"`
	VmStatus                string     `json:"vm_status"`
	AccumulatorRootHash     string     `json:"accumulator_root_hash"`
	Changes                 []Change   `json:"changes"`
	Timestamp               uint64     `json:"timestamp"`

	_ID                 string
	_Round              uint64
	_PreviousBlockVotes []bool
	_Proposer           string
}

func (t *Transaction) AsUserTransaction() *UserTransaction {
	return (*UserTransaction)(unsafe.Pointer(t))
}

type BlockMetadataTransaction struct {
	Type                     string `json:"type"`
	Hash                     string `json:"hash"`
	_Sender                  string
	_SequenceNumber          uint64
	_MaxGasAmount            uint64
	_GasUnitPrice            uint64
	_GasCurrencyCode         string
	_ExpirationTimestampSecs uint64
	_Payload                 *Payload
	_Signature               *Signature
	Events                   []Event
	Version                  uint64   `json:"version"`
	StateRootHash            string   `json:"state_root_hash"`
	EventRootHash            string   `json:"event_root_hash"`
	GasUsed                  uint64   `json:"gas_used"`
	Success                  bool     `json:"success"`
	VmStatus                 string   `json:"vm_status"`
	AccumulatorRootHash      string   `json:"accumulator_root_hash"`
	Changes                  []Change `json:"changes"`
	Timestamp                uint64   `json:"timestamp"`
	ID                       string   `json:"id"`
	Round                    uint64   `json:"round"`
	PreviousBlockVotes       []bool   `json:"previous_block_votes"`
	Proposer                 string   `json:"proposer"`
}

func (t *Transaction) AsBlockMetadataTransaction() *BlockMetadataTransaction {
	return (*BlockMetadataTransaction)(unsafe.Pointer(t))
}
