package models

import (
	"encoding/json"
	"time"
)

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (wd *Withdrawal) MarshalJSON() ([]byte, error) {
	type Alias Withdrawal
	return json.Marshal(&struct {
		*Alias
		ProcessedAt string `json:"processed_at"`
	}{
		Alias:       (*Alias)(wd),
		ProcessedAt: wd.ProcessedAt.Format(time.RFC3339),
	})
}
