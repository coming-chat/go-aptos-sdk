package aptostypes

import "unsafe"

// change 转化为特定类型 change

const (
	TypeChangeDeleteModule    = "delete_module"
	TypeChangeDeleteResource  = "delete_resource"
	TypeChangeDeleteTableItem = "delete_table_item"
	TypeChangeWriteModule     = "write_module"
	TypeChangeWriteResource   = "write_resource"
	TypeChangeWriteTableItem  = "write_table_item"
)

type DeleteModuleChange struct {
	Type         string `json:"type"`
	StateKeyHash string `json:"state_key_hash"`
	Address      string `json:"address"`
	_Resource    string
	Module       string `json:"module"`
	_Handle      string
	_Key         string
	_Value       string
	_Data        interface{}
}

type DeleteResourceChange struct {
	Type         string `json:"type"`
	StateKeyHash string `json:"state_key_hash"`
	Address      string `json:"address"`
	Resource     string `json:"resource"`
	_Module      string
	_Handle      string
	_Key         string
	_Value       string
	_Data        interface{}
}

type DeleteTableItemChange struct {
	Type         string `json:"type"`
	StateKeyHash string `json:"state_key_hash"`
	_Address     string
	_Resource    string
	_Module      string
	_Handle      string
	_Key         string
	_Value       string
	Data         interface{} `json:"data"`
}

type WriteModuleChange struct {
	Type         string `json:"type"`
	StateKeyHash string `json:"state_key_hash"`
	Address      string `json:"address"`
	_Resource    string
	_Module      string
	_Handle      string
	_Key         string
	_Value       string
	Data         interface{} `json:"data"`
}

type WriteResourceChange struct {
	Type         string `json:"type"`
	StateKeyHash string `json:"state_key_hash"`
	Address      string `json:"address"`
	_Resource    string
	_Module      string
	_Handle      string
	_Key         string
	_Value       string
	Data         interface{} `json:"data"`
}

func (wrc *WriteResourceChange) GetData() (ar AccountResource, success bool) {
	mapData, success := wrc.Data.(map[string]interface{})
	if !success {
		return
	}
	ar.Type, success = mapData["type"].(string)
	if !success {
		return
	}
	ar.Data, success = mapData["data"].(map[string]interface{})
	return
}

type WriteTableItemChange struct {
	Type         string `json:"type"`
	StateKeyHash string `json:"state_key_hash"`
	_Address     string
	_Resource    string
	_Module      string
	Handle       string `json:"handle"`
	Key          string `json:"key"`
	Value        string `json:"value"`
	_Data        interface{}
}

func (c *Change) AsDeleteModuleChange() *DeleteModuleChange {
	return (*DeleteModuleChange)(unsafe.Pointer(c))
}

func (c *Change) AsDeleteResourceChange() *DeleteResourceChange {
	return (*DeleteResourceChange)(unsafe.Pointer(c))
}

func (c *Change) AsDeleteTableItemChange() *DeleteTableItemChange {
	return (*DeleteTableItemChange)(unsafe.Pointer(c))
}

func (c *Change) AsWriteModuleChange() *WriteModuleChange {
	return (*WriteModuleChange)(unsafe.Pointer(c))
}

func (c *Change) AsWriteResourceChange() *WriteResourceChange {
	return (*WriteResourceChange)(unsafe.Pointer(c))
}

func (c *Change) AsWriteTableItemChange() *WriteTableItemChange {
	return (*WriteTableItemChange)(unsafe.Pointer(c))
}
