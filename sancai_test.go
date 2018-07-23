package fate_test

import (
	"github.com/godcong/fate"
	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2"
	"log"
	"runtime"
	"testing"
	"time"
)

func TestMakeSanCaiWuGe(t *testing.T) {
	mongo.Dial("localhost", &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	})
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	l1, l2, f1, f2 := 1, 0, 1, 0
	total := 32 * 32 * 32 * 32
	idx := 1
	max := 32
	ch := make(chan int, max)
	defer close(ch)
	for ; idx <= max; idx++ {
		go ProcessOne(ch, l1, l2, f1, f2)
		l1, l2, f1, f2 = valueLoop(l1, l2, f1, f2)
	}

	for {
		if idx > total {
			break
		}
		select {
		case <-ch:
			//log.Println("now:", v)
			go ProcessOne(ch, l1, l2, f1, f2)
			l1, l2, f1, f2 = valueLoop(l1, l2, f1, f2)
		}
	}

	//wait for all done
	time.Sleep(10 * time.Second)

}

func ProcessOne(b chan<- int, l1, l2, f1, f2 int) {
	log.Println(l1, l2, f1, f2)

	ge := fate.MakeSanCai(fate.MakeWuGe(l1, l2, f1, f2))
	err := mongo.InsertIfNotExist(mongo.C("sc"), ge)
	log.Println(err)
	b <- l1 * l2 * f1 * f2
}

func valueLoop(l1, l2, f1, f2 int) (int, int, int, int) {
	f2++

	if f2 > 32 {
		f2 = 0
		f1++
	}
	if f1 > 32 {
		f1 = 1
		l2++
	}
	if l2 > 32 {
		l2 = 0
		l1++
	}

	return l1, l2, f1, f2
}
