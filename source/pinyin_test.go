package source

import (
	"fmt"
	"testing"
)

func TestLoadPinYin(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []*PinYin
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				path: "pinyin.txt",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadPinYin(tt.args.path)
			for _, py := range got {
				fmt.Println("id", py.ID, "char", py.Char, "py", py.Pinyin)
			}
		})
	}
}
