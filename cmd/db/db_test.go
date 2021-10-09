package main

import (
	"testing"
)

func Test_writeConfig(t *testing.T) {
	type args struct {
		p  string
		db DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				p: "db.config",
				db: DB{
					From: From{
						DSN:          "",
						Host:         "localhost",
						Port:         "3306",
						User:         "root",
						Pwd:          "root",
						DBName:       "fate",
						MaxIdleCon:   0,
						MaxOpenCon:   0,
						Driver:       "mysql",
						File:         "",
						ShowSQL:      false,
						ShowExecTime: false,
					},
					To: To{
						DSN: "file:newdb.db?cache=shared&_journal=WAL&_fk=1",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeConfig(tt.args.p, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("writeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
