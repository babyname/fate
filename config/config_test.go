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

func TestGetPath(t *testing.T) {
	type args struct {
		root  string
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				root:  "E:\\",
				paths: []string{"test", "fate"},
			},
			want: `E:\test\fate`,
		},
		{
			name: "",
			args: args{
				root:  "",
				paths: []string{"test", "fate"},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPath(tt.args.root, tt.args.paths...); got != tt.want {
				t.Errorf("GetPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
