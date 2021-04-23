package forms

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/models"
	"time"
)

type AlertForm struct {
	Fingerprint string            `json:"fingerprint"`
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    *time.Time        `json:"startsAt"`
	EndsAt      *time.Time        `json:"endsAt"`
}

func (c *AlertForm) Valid(v *validation.Validation) {

}

func (c *AlertForm) ToModel() *models.Alert {
	var (
		endsAt      *time.Time
		labels      []byte
		annotations []byte
	)

	if c.EndsAt != nil && !c.EndsAt.IsZero() {
		endsAt = c.EndsAt
	}
	labels, _ = json.Marshal(c.Labels)
	annotations, _ = json.Marshal(c.Annotations)

	return &models.Alert{
		Fingerprint: c.Fingerprint,
		Instance:    c.Labels["instance"],
		AlertName:   c.Labels["alertname"],
		Severity:    c.Labels["severity"],
		Status:      c.Status,
		StartsAt:    c.StartsAt,
		EndsAt:      endsAt,
		Summary:     c.Annotations["summary"],
		Description: c.Annotations["description"],
		Labels:      string(labels),
		Annotations: string(annotations),
	}
}

type AlertsForm struct {
	Alerts []AlertForm `form:"alerts"`
}
