package model

import (
	"time"
)

// 既存テーブルに合わせてモデルを作成する場合はcolumn名は大文字小文字も合わせる必要がある。
// 合わせないとselectしたときにバインドできない。
type JobSheet struct {
	ID               string         `gorm:"primaryKey;column:ID"`
	CompleteDate     time.Time      `gorm:"column:COMPLETEDATE;type:date"`
	Content          string         `gorm:"column:CONTENT"`
	Department       string         `gorm:"column:DEPARTMENT"`
	LimitDate        time.Time      `gorm:"column:LIMITDATE;type:date"`
	OccurDateTime    time.Time      `gorm:"column:OCCURDATETIME"`
	Person           string         `gorm:"column:PERSON"`
	ResponseTime     float64        `gorm:"column:RESPONSETIME"`
	Support          string         `gorm:"column:SUPPORT"`
	Title            string         `gorm:"column:TITLE"`
	BusinessSystemID int            `gorm:"column:BUSINESSSYSTEM_ID"`
	BusinessSystem   BusinessSystem `gorm:"foreignKey:BusinessSystemID"`
	ClientID         int            `gorm:"column:CLIENT_ID"`
	Client           Client         `gorm:"foreignKey:ClientID"`
	ContactID        string         `gorm:"column:CONTACT_ID"`
	Contact          User           `gorm:"foreignKey:ContactID"`
	DealID           string         `gorm:"column:DEAL_ID"`
	Deal             User           `gorm:"foreignKey:DealID"`
	InquiryID        int            `gorm:"column:INQUIRY_ID"`
	Inquiry          Inquiry        `gorm:"foreignKey:InquiryID"`
}

func (j JobSheet) TableName() string {
	return "t_job_sheet"
}
