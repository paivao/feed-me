package types

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type NullString struct {
	sql.NullString
}

type NullTime struct {
	sql.NullTime
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid == false {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(ns.String))
	b = append(b, '"')
	b = append(b, ns.String...)
	b = append(b, '"')
	json.Marshal(b)
	return b, nil
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	sdata := string(data)
	if "null" == sdata {
		ns.Valid = false
		return nil
	}
	if sdata[0] != '"' || sdata[len(sdata)-1] != '"' {
		return errors.New("invalid string")
	}
	ns.Valid = true
	ns.String = sdata[1 : len(sdata)-1]
	return nil
}

func (ns *NullTime) MarshalJSON() ([]byte, error) {
	if ns.Valid == false {
		return []byte("null"), nil
	}
	return ns.Time.MarshalJSON()
}

func (ns *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	return ns.Time.UnmarshalJSON(data)
}
