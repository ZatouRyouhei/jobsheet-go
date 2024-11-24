package model

// 既存テーブルに合わせてモデルを作成する場合はcolumn名は大文字小文字も合わせる必要がある。
// 合わせないとselectしたときにバインドできない。
type Business struct {
	ID   int    `gorm:"primaryKey;column:ID"`
	Name string `gorm:"column:NAME"`
}

func (b Business) TableName() string {
	return "t_business"
}
