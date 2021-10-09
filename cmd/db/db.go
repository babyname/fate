package main

import (
	"encoding/json"
	"io/ioutil"
)

type DB struct {
	From From `json:"from"`
	To   To   `json:"to"`
}

func readConfig(p string) (db DB, err error) {
	filebytes, err := ioutil.ReadFile(p)
	if err != nil {
		return db, err
	}

	err = json.Unmarshal(filebytes, &db)
	if err != nil {
		return db, err
	}
	return db, nil
}

func writeConfig(p string, db DB) error {
	marshal, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, marshal, 0755)
}
