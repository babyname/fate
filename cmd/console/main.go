package main

import (
	"log"
	"runtime"
	"time"
)

func main() {

}

func initSCWG() {
	max := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(max)
	l1, l2, f1, f2 := 1, 0, 1, 0

	ch := make(chan bool, max)

	for i := 0; i < max; i++ {
		go runtineProcess(ch, l1, l2, f1, f2)
	}

	for {
		select {
		case <-ch:
			if l1 <= 32 {
				if f1 == 1 && f2 == 0 {
					log.Println("process:", l1, l2, f1, f2, time.Now())
				}
				go runtineProcess(ch, l1, l2, f1, f2)
			}
		case <-time.After(10 * time.Second):
			log.Println("time out")
			break
		}
	}

}

func runtineProcess(b chan<- bool, l1, l2, f1, f2 int) {

}
