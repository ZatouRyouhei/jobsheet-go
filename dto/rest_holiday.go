package dto

import (
	"jobsheet-go/database/model"
)

type RestHoliday struct {
	Holiday string `json:"holiday"`
	Name    string `json:"name"`
}

func NewRestHoliday(holiday model.Holiday) RestHoliday {
	var restHoliday RestHoliday
	restHoliday.Holiday = holiday.Holiday.Format("2006-01-02")
	restHoliday.Name = holiday.Name
	return restHoliday
}
