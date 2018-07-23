package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/godcong/fate/config"
	"github.com/godcong/fate/mongo"
	"github.com/globalsign/mgo"
	"github.com/godcong/fate"
)

func main() {
	args := os.Args
	if len(args) > 0 {
		switch args[0] {
		case "init":
			fallthrough
		default:
			initSCWG()
		}
	}
}

func initSCWG() {
	max := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(max)
	l1, l2, f1, f2 := 1, 0, 1, 0
	mongo.Dial(config.Default().GetStringD("mongodb.url", "localhost"),
		&mgo.Credential{
			Username: config.Default().GetStringD("mongodb.username", "root"),
			Password: config.Default().GetStringD("mongodb.password", "root"),
		})

	ch := make(chan bool, max)

	for i := 0; i < max; i++ {
		go runtineProcess(ch, l1, l2, f1, f2)
		l1, l2, f1, f2 = valuePlus(l1, l2, f1, f2)
	}

	for {
		select {
		case <-ch:
			if l1 <= 32 {
				if f1 == 1 && f2 == 0 {
					log.Println("process:", l1, l2, f1, f2, time.Now())
				}
				go runtineProcess(ch, l1, l2, f1, f2)
				l1, l2, f1, f2 = valuePlus(l1, l2, f1, f2)
			}
		case <-time.After(10 * time.Second):
			log.Println("time out")
			break
		}
	}

}

func runtineProcess(b chan<- bool, l1, l2, f1, f2 int) {
	wuGe := fate.MakeWuGe(l1, l2, f1, f2)
	sanCai := fate.MakeSanCai(wuGe)
	mongo.InsertIfNotExist(mongo.C("sancai"), sanCai)
	b <- true
}

func valuePlus(l1, l2, f1, f2 int) (int, int, int, int) {
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
