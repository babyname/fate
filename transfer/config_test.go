package transfer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/babyname/fate/config"
)

func Test_writeConfig(t *testing.T) {
	type args struct {
		p  string
		db DatabaseConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				p: "transfer.cfg",
				db: DatabaseConfig{
					Tables: []string{"Character", "WuGeLucky", "WuXing"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteTransferConfig(tt.args.p, &tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("writeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadTransferConfig(t *testing.T) {
	want := &DatabaseConfig{
		Tables: []string{"Character", "WuGeLucky", "WuXing"},
	}
	err := WriteTransferConfig("transfer.cfg", want)
	if err != nil {
		t.Fatalf("WriteTransferConfig() error = %v", err)
	}
	got, err := ReadTransferConfig("transfer.cfg")
	if err != nil {
		t.Fatalf("ReadTransferConfig() error = %v", err)
	}
	if !reflect.DeepEqual(got.Tables, want.Tables) {
		marshal, _ := json.Marshal(got)
		t.Errorf("ReadTransferConfig() got = %s, want Tables=%v", string(marshal), want.Tables)
	}
}

func TestDatabaseConfigUnmarshal(t *testing.T) {
	cfg := config.DefaultConfig()
	raw, _ := json.Marshal(cfg)

	var db DatabaseConfig
	db.SourceRaw = raw
	db.TargetRaw = raw
	db.Tables = []string{"Character"}

	err := json.Unmarshal(db.SourceRaw, &db.Source)
	if err != nil {
		t.Fatalf("unmarshal source: %v", err)
	}
	if db.Source.Database.Driver != "sqlite3" {
		t.Errorf("Source.Database.Driver = %v, want sqlite3", db.Source.Database.Driver)
	}
}
