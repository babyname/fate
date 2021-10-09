package main

import (
	"reflect"
	"testing"
)

func Test_runTransfer(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runTransfer(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("runTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readConfig(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantDb  DB
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				p: "db.config",
			},
			wantDb:  DB{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDb, err := readConfig(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDb, tt.wantDb) {
				t.Errorf("readConfig() gotDb = %v, want %v", gotDb, tt.wantDb)
			}
		})
	}
}
