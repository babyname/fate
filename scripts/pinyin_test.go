package scripts

import (
	"testing"

	"github.com/babyname/fate/log"
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
		{
			name: "",
			args: args{
				path: "zdic.txt",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadPinYin(tt.args.path, func(yin *PinYin) bool {
				if yin.ID < 10000 {
					log.Logger("scripts").Info("log string", "id", yin.ID, "pinyin", yin.Pinyin, "char", yin.Char)
				}
				return true
			})

		})
	}
}
