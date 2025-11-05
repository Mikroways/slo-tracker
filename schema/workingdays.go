package schema

type WorkingDaySchedule struct {
	Weekday   int    `json:"weekday"`
	OpenHour  string `json:"open_hour"`
	CloseHour string `json:"close_hour"`
}

type StoreWorkingSchedule struct {
	ID        uint   `json:"id omitempty" gorm:"primaryKey"`
	SLOID     uint   `json:"slo_id" gorm:"index"`
	Weekday   int    `json:"weekday"`
	OpenHour  string `json:"open_hour" gorm:"type:time"`
	CloseHour string `json:"close_hour" gorm:"type:time"`
}
