package aptostypes

const (
	EntryFunctionPayload = "entry_function_payload"
	ScriptPayload        = "script_payload"
	ModuleBundlePayload  = "module_bundle_payload"
	WriteSetPayload      = "write_set_payload"
)

//go:generate go run github.com/fjl/gencodec -type RestError -field-override restErrorMarshaling -out gen_rest_error_json.go

type RestError struct {
	Code               int    `json:"code"`
	Message            string `json:"message"`
	AptosLedgerVersion uint64 `json:"aptos_ledger_version"`
}

func (e *RestError) Error() string {
	return e.Message
}

type restErrorMarshaling struct {
	AptosLedgerVersion jsonUint64
}

//go:generate go run github.com/fjl/gencodec -type LedgerInfo -field-override ledgerInfoMarshaling -out gen_ledger_json.go

// LedgerInfo represents chain current ledger info
type LedgerInfo struct {
	ChainId             int    `json:"chain_id"`
	LedgerVersion       uint64 `json:"ledger_version" gencodec:"required"`
	LedgerTimestamp     uint64 `json:"ledger_timestamp" gencodec:"required"`
	BlockHeight         uint64 `json:"block_height" gencodec:"required"`
	Epoch               uint64 `json:"epoch"`
	NodeRole            string `json:"node_role"`
	OldestBlockHeight   uint64 `json:"oldest_block_height"`
	OldestLedgerVersion uint64 `json:"oldest_ledger_version"`
}

type ledgerInfoMarshaling struct {
	Epoch               jsonUint64
	LedgerVersion       jsonUint64
	LedgerTimestamp     jsonUint64
	BlockHeight         jsonUint64
	OldestBlockHeight   jsonUint64
	OldestLedgerVersion jsonUint64
}

//go:generate go run github.com/fjl/gencodec -type Block -field-override blockMarshaling -out gen_block_json.go

type Block struct {
	BlockHeight    uint64        `json:"block_height"`
	BlockHash      string        `json:"block_hash"`
	BlockTimestamp uint64        `json:"block_timestamp"`
	FirstVersion   uint64        `json:"first_version"`
	LastVersion    uint64        `json:"last_version"`
	Transactions   []Transaction `json:"transactions"`
}

type blockMarshaling struct {
	BlockHeight    jsonUint64
	BlockTimestamp jsonUint64
	FirstVersion   jsonUint64
	LastVersion    jsonUint64
}

//go:generate go run github.com/fjl/gencodec -type Transaction -field-override transactionMarshaling -out gen_transaction_json.go

// Transaction
type Transaction struct {
	Type                    string     `json:"type"`
	Hash                    string     `json:"hash"`                      // PendingTransaction|GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	Sender                  string     `json:"sender"`                    // PendingTransaction|UserTransaction
	SequenceNumber          uint64     `json:"sequence_number"`           // PendingTransaction|UserTransaction
	MaxGasAmount            uint64     `json:"max_gas_amount"`            // PendingTransaction|UserTransaction
	GasUnitPrice            uint64     `json:"gas_unit_price"`            // PendingTransaction|UserTransaction
	GasCurrencyCode         string     `json:"gas_currency_code"`         // PendingTransaction|UserTransaction
	ExpirationTimestampSecs uint64     `json:"expiration_timestamp_secs"` // PendingTransaction|UserTransaction
	Payload                 *Payload   `json:"payload"`                   // PendingTransaction|GenesisTransaction|UserTransaction
	Signature               *Signature `json:"signature"`                 // PendingTransaction|UserTransaction

	Events              []Event  `json:"events"`                // GenesisTransaction|UserTransaction|BlockMetadataTransaction
	Version             uint64   `json:"version"`               // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	StateRootHash       string   `json:"state_root_hash"`       // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	EventRootHash       string   `json:"event_root_hash"`       // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	GasUsed             uint64   `json:"gas_used"`              // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	Success             bool     `json:"success"`               // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	VmStatus            string   `json:"vm_status"`             // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	AccumulatorRootHash string   `json:"accumulator_root_hash"` // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
	Changes             []Change `json:"changes"`               // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction

	Timestamp uint64 `json:"timestamp"` // UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction

	ID                 string `json:"id"`                   //BlockMetadataTransaction
	Round              uint64 `json:"round"`                //BlockMetadataTransaction
	PreviousBlockVotes []bool `json:"previous_block_votes"` //BlockMetadataTransaction
	Proposer           string `json:"proposer"`             //BlockMetadataTransaction
}

type transactionMarshaling struct {
	SequenceNumber          jsonUint64
	MaxGasAmount            jsonUint64
	GasUnitPrice            jsonUint64
	ExpirationTimestampSecs jsonUint64
	GasUsed                 jsonUint64
	Version                 jsonUint64
	Round                   jsonUint64
	Timestamp               jsonUint64
}

type Payload struct {
	Type          string        `json:"type"`           // ScriptPayload|ScriptFunctionPayload|ModuleBundlePayload|WriteSetPayload
	TypeArguments []string      `json:"type_arguments"` // ScriptPayload|ScriptFunctionPayload
	Arguments     []interface{} `json:"arguments"`      // ScriptPayload|ScriptFunctionPayload

	Modules []MoveModule `json:"modules"` // ModuleBundlePayload

	Code *MoveScript `json:"code"` // ScriptPayload

	Function string `json:"function"` // ScriptFunctionPayload

	WriteSet *WriteSet `json:"write_set"` // WriteSetPayload
}

type (
	WriteSet struct {
		Type      string  `json:"type"`       // script_write_set|direct_write_set
		ExecuteAs string  `json:"execute_as"` // script_write_set
		Script    *Script `json:"script"`     // script_write_set

		Changes []Change `json:"changes"` // direct_write_set
		Events  []Event  `json:"events"`  // direct_write_set
	}

	Script struct {
		Code          *MoveScript   `json:"code"`
		TypeArguments []string      `json:"type_arguments"`
		Arguments     []interface{} `json:"arguments"`
	}
)

type (
	MoveModule struct {
		ByteCode string         `json:"bytecode"`
		Abi      *MoveModuleAbi `json:"abi"`
	}

	MoveModuleAbi struct {
		Address          string         `json:"address"`
		Name             string         `json:"name"`
		Friend           []string       `json:"friend"`
		ExposedFunctions []MoveFunction `json:"exposed_functions"`
		Structs          []MoveStruct   `json:"structs"`
	}

	MoveStruct struct {
		Name              string
		IsNative          bool
		Abilities         []string          `json:"abilities"` // copy|drop|store|key
		GenericTypeParams []interface{}     `json:"generic_type_params"`
		Fields            []MoveStructField `json:"fields"`
	}

	MoveStructField struct {
		Name string
		Type string // eg.0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>, match ^(bool|u8|u64|u128|address|signer|vector<.+>|0x[0-9a-zA-Z:_<, >]+|^&(mut )?.+$|T\d+)$
	}
)

type MoveScript struct {
	Bytecode string        `json:"bytecode"`
	Abi      *MoveFunction `json:"abi"`
}

type MoveFunction struct {
	Name              string        `json:"name"`
	Visibility        string        `json:"visibility"` // public|script|friend
	GenericTypeParams []interface{} `json:"generic_type_params"`
	Params            []string      `json:"params"`
	Return            []string      `json:"return"`
}

type (
	Change struct {
		Type         string `json:"type"`           // delete_module|delete_resource|delete_table_item|write_module|write_resource|write_table_item
		StateKeyHash string `json:"state_key_hash"` // delete_module|delete_resource|delete_table_item|write_module
		Address      string `json:"address"`        // delete_module|delete_resource|write_module
		Resource     string `json:"resource"`       // delete_resource
		Module       string `json:"module"`         // delete_module

		Handle string `json:"handle"` //write_table_item
		Key    string `json:"key"`    // write_table_item
		Value  string `json:"value"`  //write_table_item

		Data interface{} `json:"data"` // delete_table_item(DeleteTableItem)|write_module(MoveModule)|write_resource(AccountResource)
	}

	DeleteTableItem struct {
		Handle string `json:"handle"`
		Key    string `json:"key"`
	}
)

//go:generate go run github.com/fjl/gencodec -type Event -field-override eventMarshaling -out gen_event_json.go

type Event struct {
	Key            string      `json:"key"`
	SequenceNumber uint64      `json:"sequence_number"`
	Type           string      `json:"type"` // eg. 0x1::aptos_coin::AptosCoin, match ^(bool|u8|u64|u128|address|signer|vector<.+>|0x[0-9a-zA-Z:_<, >]+)$
	Data           interface{} `json:"data"`
	Version        uint64      `json:"version"` // only getEvents will return version
}

type eventMarshaling struct {
	SequenceNumber jsonUint64
	Version        jsonUint64
}

type Signature struct {
	Type      string `json:"type"`       // ed25519_signature|multi_ed25519_signature|multi_agent_signature
	PublicKey string `json:"public_key"` // ed25519_signature
	Signature string `json:"signature"`  // ed25519_signature

	PublicKeys []string `json:"public_keys"` // multi_ed25519_signature
	Signatures []string `json:"signatures"`  // multi_ed25519_signature
	Threshold  int      `json:"threshold"`   // multi_ed25519_signature
	Bitmap     string   `json:"bitmap"`      // multi_ed25519_signature

	Sender                   *Signature  `json:"Sender"`                     // multi_agent_signature
	SecondarySignerAddresses []string    `json:"secondary_signer_addresses"` // multi_agent_signature
	SecondarySigners         []Signature `json:"secondary_signers"`          // multi_agent_signature
}

//go:generate go run github.com/fjl/gencodec -type AccountCoreData -field-override AccountCoreDataMarshaling -out gen_account_core_json.go

type AccountCoreData struct {
	SequenceNumber    uint64 `json:"sequence_number"`
	AuthenticationKey string `json:"authentication_key"`
}

type AccountCoreDataMarshaling struct {
	SequenceNumber jsonUint64
}

type AccountResource struct {
	Type string                 `json:"type"` // match ^0x[0-9a-zA-Z:_<>]+$
	Data map[string]interface{} `json:"data"`
}
