package source

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []PolyPhone
	}{
		{
			name: "",
			args: args{
				path: "char_detail.json",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadCharJSON(tt.args.path, func(phone PolyPhone) bool {
				return false
			})
		})
	}
}

func TestCharCode(t *testing.T) {
	fmt.Println(int32(rune('靐')))
	fmt.Printf("%0x", int32(rune('靐')))
}
