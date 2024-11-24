package model

import (
	"time"
)

// 既存テーブルに合わせてモデルを作成する場合はcolumn名は大文字小文字も合わせる必要がある。
// 合わせないとselectしたときにバインドできない。
type Holiday struct {
	Holiday time.Time `gorm:"column:HOLIDAY;type:date"`
	Name    string    `gorm:"column:NAME"`
}

func (b Holiday) TableName() string {
	return "t_holiday"
}
