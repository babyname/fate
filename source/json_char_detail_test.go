package source

import "testing"

func TestLoadCharDetailJSON(t *testing.T) {
	type args struct {
		path string
		hook func(ch Character) bool
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
			if err := LoadCharDetailJSON(tt.args.path, tt.args.hook); (err != nil) != tt.wantErr {
				t.Errorf("LoadCharDetailJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
