package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		wantC *Config
	}{
		{
			name:  "",
			args:  args{},
			wantC: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotC := LoadConfig(tt.args.path); !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("LoadConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	type args struct {
		path   string
		config Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				path:   "test.mysql.config",
				config: DefaultConfig(),
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				path:   "test.sqlite.config",
				config: DefaultSqliteConfig(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.config.Save(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("SaveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
