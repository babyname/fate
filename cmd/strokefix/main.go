package main

import "github.com/godcong/fate"

func main() {
	var e error
	db := fate.InitMysql("127.0.0.1:3306", "root", "111111")
	e = db.Sync2(fate.Character{})
	if e != nil {
		panic(e)
	}
	e = UpdateFix(db)
	if e != nil {
		panic(e)
	}
}
