package model

import (
	"time"
)

type JobSheetEx struct {
	ID               string    `gorm:"primaryKey;column:id"`
	CompleteDate     time.Time `gorm:"column:completedate;type:date"`
	Content          string    `gorm:"column:content"`
	Department       string    `gorm:"column:department"`
	LimitDate        time.Time `gorm:"column:limitdate;type:date"`
	OccurDateTime    time.Time `gorm:"column:occurdatetime"`
	Person           string    `gorm:"column:person"`
	ResponseTime     float64   `gorm:"column:responsetime"`
	Support          string    `gorm:"column:support"`
	Title            string    `gorm:"column:title"`
	BusinessSystemID int       `gorm:"column:businesssystem_id"`
	ClientID         int       `gorm:"column:client_id"`
	ContactID        string    `gorm:"column:contact_id"`
	DealID           string    `gorm:"column:deal_id"`
	InquiryID        int       `gorm:"column:inquiry_id"`
}

func (j JobSheetEx) TableName() string {
	return "t_job_sheet_ex"
}
