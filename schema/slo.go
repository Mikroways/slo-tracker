package schema

import (
	"time"

	"gorm.io/gorm"
)

// SLO stores the SLO response payload
type SLO struct {
	gorm.Model

	ID        uint       `json:"id,omitempty" sql:"primary_key"`
	SLOName   string     `json:"slo_name" gorm:"unique;not null"`
	OpenHour  string     `json:"open_hour" gorm:"type:time"`
	CloseHour string     `json:"close_hour" gorm:"type:time"`
	TargetSLO float32    `json:"target_slo"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" sql:"default:current_timestamp"`
}

type SLOResponse struct {
	ID                 uint    `json:"id,omitempty"`
	SLOName            string  `json:"slo_name"`
	OpenHour           string  `json:"open_hour"`
	CloseHour          string  `json:"close_hour"`
	TargetSLO          float32 `json:"target_slo"`
	CurrentSLO         float32 `json:"current_slo"`
	RemainingErrBudget float32 `json:"remaining_err_budget"`
}
