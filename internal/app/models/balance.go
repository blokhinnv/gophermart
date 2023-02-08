package models

import (
	"database/sql"
	"encoding/json"
)

type Balance struct {
	Current   sql.NullFloat64 `json:"current"`
	Withdrawn sql.NullFloat64 `json:"withdrawn"`
}

func (b *Balance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		Current:   b.Current.Float64,
		Withdrawn: b.Withdrawn.Float64,
	})
}
