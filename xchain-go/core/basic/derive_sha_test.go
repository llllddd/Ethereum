package basic

import (
	"testing"
)

func TestDeriveSha(t *testing.T) {

	txs := Mocktxs()
	// txlist:=
	var txempty []*Transaction
	tests := []struct {
		name string
		txs  []*Transaction
		want string
	}{
		{"两笔交易", txs, "0xe543b7b837144e7be43fc17fa86c6053cb4604cf1cc4fb3f776e482a64818262"},
		// {"交易列表", txs, ""},
		{"空交易", txempty, "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveSha(tt.txs)
			if got.String() != tt.want {
				t.Error("构造的hash不相符，got= ", got.String(), ",want= ", tt.want)
			}
		})

	}
}
