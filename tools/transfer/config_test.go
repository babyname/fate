package transfer

import (
	"reflect"
	"testing"

	"github.com/tikafog/jsongs"

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
					Source: *config.DefaultMysqlConfig(),
					Target: *config.DefaultConfig(),
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
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    *DatabaseConfig
		wantErr bool
	}{
		{
			name: "",
			args: args{
				p: "transfer.cfg",
			},
			want: &DatabaseConfig{
				Tables: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadTransferConfig(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTransferConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadTransferConfig() got = %v, want %v", got, tt.want)
			}
			marshal, err := jsongs.Marshal(got)
			if err != nil {
				t.Errorf("ReadTransferConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(string(marshal))
		})
	}
}
