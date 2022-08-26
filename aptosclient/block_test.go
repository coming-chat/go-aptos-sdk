package aptosclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func TestRestClient_GetBlockByHeight(t *testing.T) {
	url := "https://fullnode.devnet.aptoslabs.com"
	client, err := Dial(context.Background(), url)
	if err != nil {
		panic(err)
	}
	type fields struct {
		chainId int
		c       *http.Client
		rpcUrl  string
		version string
	}
	type args struct {
		height            string
		with_transactions bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantBlock *aptostypes.Block
		wantErr   bool
	}{
		{
			name: "test with tx",
			fields: fields{
				chainId: 1,
				c:       client.c,
				rpcUrl:  url,
				version: VERSION1,
			},
			args: args{
				height:            "10",
				with_transactions: false,
			},
			wantBlock: &aptostypes.Block{BlockHeight: 10},
			wantErr:   false,
		},
		{
			name: "test without tx",
			fields: fields{
				chainId: 1,
				c:       client.c,
				rpcUrl:  url,
				version: VERSION1,
			},
			args: args{
				height:            "11",
				with_transactions: true,
			},
			wantBlock: &aptostypes.Block{BlockHeight: 11, Transactions: []aptostypes.Transaction{{}}},
			wantErr:   false,
		},
		{
			name: "test error",
			fields: fields{
				chainId: 1,
				c:       client.c,
				rpcUrl:  url,
				version: VERSION1,
			},
			args: args{
				height:            "-1",
				with_transactions: true,
			},
			wantBlock: &aptostypes.Block{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RestClient{
				chainId: tt.fields.chainId,
				c:       tt.fields.c,
				rpcUrl:  tt.fields.rpcUrl,
				version: tt.fields.version,
			}
			gotBlock, err := c.GetBlockByHeight(tt.args.height, tt.args.with_transactions)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestClient.GetBlockByHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBlock.BlockHeight != tt.wantBlock.BlockHeight {
				t.Errorf("RestClient.GetBlockByHeight() = %v, want %v", gotBlock, tt.wantBlock)
			}
			if len(tt.wantBlock.Transactions) > 0 && len(gotBlock.Transactions) == 0 {
				t.Errorf("RestClient.GetBlockByHeight() transaction len fail 01")
			}
			if len(tt.wantBlock.Transactions) == 0 && len(gotBlock.Transactions) != 0 {
				t.Errorf("RestClient.GetBlockByHeight() transaction len fail 02")
			}
		})
	}
}
