package ent

//#go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/lock ./schema
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/lock --template ./template --template glob=./template/*.tmpl ./schema
