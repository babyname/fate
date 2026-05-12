package wuge_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/internal/wuge"
)

func TestWuGe_WaiGe(t *testing.T) {
	l1, l2, f1, f2 := 1, 1, 1, 1
	for i := 0; i < 80000; i++ {
		if f2 >= repository.WuGeLuckyMax {
			f1++
			f2 = 1
		}
		if f1 >= repository.WuGeLuckyMax {
			l2++
			f1 = 1
		}
		if l2 >= repository.WuGeLuckyMax {
			l1++
			l2 = 1
		}
		wg := wuge.CalcWuGe(l1, l2, f1, f2)
		sum := l1 + l2 + f1 + f2
		if wg.ZongGe() != sum {
			log.Println(wg.ZongGe() == sum, l1, l2, f1, f2, wg.ZongGe())
		}
		f2++
	}
}

func TestWuGeID(t *testing.T) {
	type args struct {
		l1 int
		l2 int
		f1 int
		f2 int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				l1: 0x2,
				l2: 0x4,
				f1: 0x8,
				f2: 0x16,
			},
			want: 0x02040816,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repository.WuGeLuckyID(tt.args.l1, tt.args.l2, tt.args.f1, tt.args.f2)
			fmt.Printf("got %x\n", got)
			if got != tt.want {
				t.Errorf("WuGeID() = %v, nowant %v", got, tt.want)
			}
		})
	}
}
