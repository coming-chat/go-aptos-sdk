package aptostypes

import (
	"reflect"
	"testing"
)

func TestChange_AsDeleteModuleChange(t *testing.T) {
	type fields struct {
		Type         string
		StateKeyHash string
		Address      string
		Resource     string
		Module       string
		Handle       string
		Key          string
		Value        string
		Data         interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *DeleteModuleChange
	}{
		{
			name: "test",
			fields: fields{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				Resource:     "resourcesxxxxx",
				Module:       "modulexx",
				Handle:       "handlexxxx",
				Key:          "xxxkey",
				Value:        "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
			want: &DeleteModuleChange{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				_Resource:    "resourcesxxxxx",
				Module:       "modulexx",
				_Handle:      "handlexxxx",
				_Key:         "xxxkey",
				_Value:       "valuexx",
				_Data:        map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Change{
				Type:         tt.fields.Type,
				StateKeyHash: tt.fields.StateKeyHash,
				Address:      tt.fields.Address,
				Resource:     tt.fields.Resource,
				Module:       tt.fields.Module,
				Handle:       tt.fields.Handle,
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				Data:         tt.fields.Data,
			}
			if got := c.AsDeleteModuleChange(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.AsDeleteModuleChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChange_AsDeleteResourceChange(t *testing.T) {
	type fields struct {
		Type         string
		StateKeyHash string
		Address      string
		Resource     string
		Module       string
		Handle       string
		Key          string
		Value        string
		Data         interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *DeleteResourceChange
	}{
		{
			name: "test",
			fields: fields{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				Resource:     "resourcesxxxxx",
				Module:       "modulexx",
				Handle:       "handlexxxx",
				Key:          "xxxkey",
				Value:        "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
			want: &DeleteResourceChange{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				Resource:     "resourcesxxxxx",
				_Module:      "modulexx",
				_Handle:      "handlexxxx",
				_Key:         "xxxkey",
				_Value:       "valuexx",
				_Data:        map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Change{
				Type:         tt.fields.Type,
				StateKeyHash: tt.fields.StateKeyHash,
				Address:      tt.fields.Address,
				Resource:     tt.fields.Resource,
				Module:       tt.fields.Module,
				Handle:       tt.fields.Handle,
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				Data:         tt.fields.Data,
			}
			if got := c.AsDeleteResourceChange(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.AsDeleteResourceChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChange_AsDeleteTableItemChange(t *testing.T) {
	type fields struct {
		Type         string
		StateKeyHash string
		Address      string
		Resource     string
		Module       string
		Handle       string
		Key          string
		Value        string
		Data         interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *DeleteTableItemChange
	}{
		{
			name: "test",
			fields: fields{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				Resource:     "resourcesxxxxx",
				Module:       "modulexx",
				Handle:       "handlexxxx",
				Key:          "xxxkey",
				Value:        "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
			want: &DeleteTableItemChange{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				_Address:     "0xa123",
				_Resource:    "resourcesxxxxx",
				_Module:      "modulexx",
				_Handle:      "handlexxxx",
				_Key:         "xxxkey",
				_Value:       "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Change{
				Type:         tt.fields.Type,
				StateKeyHash: tt.fields.StateKeyHash,
				Address:      tt.fields.Address,
				Resource:     tt.fields.Resource,
				Module:       tt.fields.Module,
				Handle:       tt.fields.Handle,
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				Data:         tt.fields.Data,
			}
			if got := c.AsDeleteTableItemChange(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.AsDeleteTableItemChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChange_AsWriteModuleChange(t *testing.T) {
	type fields struct {
		Type         string
		StateKeyHash string
		Address      string
		Resource     string
		Module       string
		Handle       string
		Key          string
		Value        string
		Data         interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *WriteModuleChange
	}{
		{
			name: "test",
			fields: fields{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				Resource:     "resourcesxxxxx",
				Module:       "modulexx",
				Handle:       "handlexxxx",
				Key:          "xxxkey",
				Value:        "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
			want: &WriteModuleChange{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				_Resource:    "resourcesxxxxx",
				_Module:      "modulexx",
				_Handle:      "handlexxxx",
				_Key:         "xxxkey",
				_Value:       "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Change{
				Type:         tt.fields.Type,
				StateKeyHash: tt.fields.StateKeyHash,
				Address:      tt.fields.Address,
				Resource:     tt.fields.Resource,
				Module:       tt.fields.Module,
				Handle:       tt.fields.Handle,
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				Data:         tt.fields.Data,
			}
			if got := c.AsWriteModuleChange(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.AsWriteModuleChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChange_AsWriteResourceChange(t *testing.T) {
	type fields struct {
		Type         string
		StateKeyHash string
		Address      string
		Resource     string
		Module       string
		Handle       string
		Key          string
		Value        string
		Data         interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *WriteResourceChange
	}{
		{
			name: "test",
			fields: fields{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				Resource:     "resourcesxxxxx",
				Module:       "modulexx",
				Handle:       "handlexxxx",
				Key:          "xxxkey",
				Value:        "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
			want: &WriteResourceChange{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				_Resource:    "resourcesxxxxx",
				_Module:      "modulexx",
				_Handle:      "handlexxxx",
				_Key:         "xxxkey",
				_Value:       "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Change{
				Type:         tt.fields.Type,
				StateKeyHash: tt.fields.StateKeyHash,
				Address:      tt.fields.Address,
				Resource:     tt.fields.Resource,
				Module:       tt.fields.Module,
				Handle:       tt.fields.Handle,
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				Data:         tt.fields.Data,
			}
			if got := c.AsWriteResourceChange(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.AsWriteResourceChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChange_AsWriteTableItemChange(t *testing.T) {
	type fields struct {
		Type         string
		StateKeyHash string
		Address      string
		Resource     string
		Module       string
		Handle       string
		Key          string
		Value        string
		Data         interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *WriteTableItemChange
	}{
		{
			name: "test",
			fields: fields{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				Address:      "0xa123",
				Resource:     "resourcesxxxxx",
				Module:       "modulexx",
				Handle:       "handlexxxx",
				Key:          "xxxkey",
				Value:        "valuexx",
				Data:         map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
			want: &WriteTableItemChange{
				Type:         TypeChangeDeleteModule,
				StateKeyHash: "keyhash123123xxx",
				_Address:     "0xa123",
				_Resource:    "resourcesxxxxx",
				_Module:      "modulexx",
				Handle:       "handlexxxx",
				Key:          "xxxkey",
				Value:        "valuexx",
				_Data:        map[string]interface{}{"type": "xxx", "data": "xxx"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Change{
				Type:         tt.fields.Type,
				StateKeyHash: tt.fields.StateKeyHash,
				Address:      tt.fields.Address,
				Resource:     tt.fields.Resource,
				Module:       tt.fields.Module,
				Handle:       tt.fields.Handle,
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				Data:         tt.fields.Data,
			}
			if got := c.AsWriteTableItemChange(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.AsWriteTableItemChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteResourceChange_GetData(t *testing.T) {
	type fields struct {
		Type         string
		StateKeyHash string
		Address      string
		_Resource    string
		_Module      string
		_Handle      string
		_Key         string
		_Value       string
		Data         interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		wantAr      AccountResource
		wantSuccess bool
	}{
		{
			name: "success case",
			fields: fields{
				Data: map[string]interface{}{
					"type": "xx",
					"data": map[string]interface{}{"a": 1},
				},
			},
			wantAr: AccountResource{
				Type: "xx",
				Data: map[string]interface{}{"a": 1},
			},
			wantSuccess: true,
		},
		{
			name: "fail case",
			fields: fields{
				Data: map[string]interface{}{
					"type": "xx",
				},
			},
			wantAr: AccountResource{
				Type: "xx",
			},
			wantSuccess: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrc := &WriteResourceChange{
				Type:         tt.fields.Type,
				StateKeyHash: tt.fields.StateKeyHash,
				Address:      tt.fields.Address,
				_Resource:    tt.fields._Resource,
				_Module:      tt.fields._Module,
				_Handle:      tt.fields._Handle,
				_Key:         tt.fields._Key,
				_Value:       tt.fields._Value,
				Data:         tt.fields.Data,
			}
			gotAr, gotSuccess := wrc.GetData()
			if !reflect.DeepEqual(gotAr, tt.wantAr) {
				t.Errorf("WriteResourceChange.GetData() gotAr = %v, want %v", gotAr, tt.wantAr)
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("WriteResourceChange.GetData() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}
