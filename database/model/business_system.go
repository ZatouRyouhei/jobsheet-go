package model

// 既存テーブルに合わせてモデルを作成する場合はcolumn名は大文字小文字も合わせる必要がある。
// 合わせないとselectしたときにバインドできない。
type BusinessSystem struct {
	ID         int      `gorm:"primaryKey;column:ID"`
	Name       string   `gorm:"column:NAME"`
	BusinessID int      `gorm:"column:BUSINESS_ID"`
	Business   Business `gorm:"foreignKey:BusinessID"`
}

func (b BusinessSystem) TableName() string {
	return "t_business_system"
}
