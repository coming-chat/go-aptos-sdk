package aptostypes

const (
	PendingTransaction            = "pending_transaction"
	GenesisTransaction            = "genesis_transaction"
	UserTransaction               = "user_transaction"
	BlockMetadataTransaction      = "block_metadata_transaction"
	StateCheckoutpointTransaction = "state_checkpoint_transaction"

	ScriptFunctionPayload = "script_function_payload"
	ScriptPayload         = "script_payload"
	ModuleBundlePayload   = "module_bundle_payload"
	WriteSetPayload       = "write_set_payload"
)

//go:generate go run github.com/fjl/gencodec -type LedgerInfo -field-override ledgerInfoMarshaling -out gen_ledger_json.go

// LedgerInfo represents chain current ledger info
type LedgerInfo struct {
	ChainId         int    `json:"chain_id"`
	LedgerVersion   uint64 `json:"ledger_version" gencodec:"required"`
	LedgerTimestamp uint64 `json:"ledger_timestamp" gencodec:"required"`

	Epoch    int    `json:"epoch"`
	NodeRole string `json:"node_role"`
}

type ledgerInfoMarshaling struct {
	LedgerVersion   Uint64
	LedgerTimestamp Uint64
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

	Events              []Event  `json:"events"`                // GenesisTransaction|UserTransaction
	Version             Uint64   `json:"version"`               // GenesisTransaction|UserTransaction|BlockMetadataTransaction|StateCheckoutpointTransaction
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
	SequenceNumber          Uint64
	MaxGasAmount            Uint64
	GasUnitPrice            Uint64
	ExpirationTimestampSecs Uint64
	GasUsed                 Uint64
	Version                 Uint64
	Round                   Uint64
	Timestamp               Uint64
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
		Structs          *MoveStruct    `json:"structs"`
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
		StateKeyHash string `json:"state_key_hash"` // delete_module|delete_resource
		Address      string `json:"address"`        // delete_module|delete_resource
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

	AccountResource struct {
		Type string      `json:"type"` // move type, match: ^0x[0-9a-zA-Z:_<>]+$
		Data interface{} `json:"data"`
	}
)

type Event struct {
	Key            string      `json:"key"`
	SequenceNumber string      `json:"sequence_number"`
	Type           string      `json:"type"` // eg. 0x1::aptos_coin::AptosCoin, match ^(bool|u8|u64|u128|address|signer|vector<.+>|0x[0-9a-zA-Z:_<, >]+)$
	Data           interface{} `json:"data"`
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
