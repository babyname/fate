package source

import (
	"testing"
)

func init() {
	//cli := database.New(config.DefaultConfig().Database)
	//client, err := cli.Client()
	//if err != nil {
	//	panic(err)
	//}
}

func TestLoadWord(t *testing.T) {
	type args struct {
		path string
		hook func(w Word) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				path: "",
				hook: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadWord(tt.args.path, tt.args.hook); (err != nil) != tt.wantErr {
				t.Errorf("LoadWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
