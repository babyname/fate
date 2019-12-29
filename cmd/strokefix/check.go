package main

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"io/ioutil"
)

var checks map[int]string

func CheckLoader(s string) error {
	bytes, e := ioutil.ReadFile(s)
	if e != nil {
		return e
	}

	e = json.Unmarshal(bytes, &checks)
	if e != nil {
		return e
	}

	if len(checks) == 0 {
		return xerrors.New("no data")
	}

	return nil
}
