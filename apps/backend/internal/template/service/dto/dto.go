package dto

import (
	"encoding/json"
	"time"

	"problum/internal/api"
	"problum/internal/model"
)

type Template struct {
	ID        int
	ProblemID int
	Language  string
	Code      string
	Metadata  json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToDTO(template *model.Template) *Template {
	return &Template{
		ID:        template.ID,
		ProblemID: template.ProblemID,
		Language:  template.Language,
		Code:      template.Code,
		Metadata:  template.Metadata,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}
}

func ToAPI(template *Template) api.TemplateGetResponse {
	if template == nil {
		return api.TemplateGetResponse{}
	}

	return api.TemplateGetResponse{
		ID:        template.ID,
		ProblemID: template.ProblemID,
		Language:  template.Language,
		Code:      template.Code,
		Metadata:  template.Metadata,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}
}
