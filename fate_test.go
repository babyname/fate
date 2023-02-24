package fate

import (
	"testing"

	"github.com/babyname/fate/config"

	_ "github.com/mattn/go-sqlite3"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		nowant  Fate
		wantErr bool
	}{
		{
			name: "",
			args: args{
				cfg: config.DefaultConfig(),
			},
			nowant:  nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.nowant {
				t.Errorf("New() got = %v, nowant %v", got, tt.nowant)
			}
		})
	}
}
