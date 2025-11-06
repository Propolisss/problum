package api

import (
	"encoding/json"
	"time"
)

type TemplateGetResponse struct {
	ID        int             `json:"id"`
	ProblemID int             `json:"problem_id"`
	Language  string          `json:"language"`
	Code      string          `json:"code"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
