package scripts

import "testing"

func TestLoadKangXiChar(t *testing.T) {
	type args struct {
		path string
		hook func(kx KangXi) bool
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
			if err := LoadKangXiChar(tt.args.path, tt.args.hook); (err != nil) != tt.wantErr {
				t.Errorf("LoadKangXiChar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
