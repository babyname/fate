package transfer

import (
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
				p: "db.config",
				db: DatabaseConfig{
					Source: config.DefaultConfig(),
					Target: config.DefaultSqliteConfig(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeConfig(tt.args.p, &tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("writeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
