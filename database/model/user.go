package model

// 既存テーブルに合わせてモデルを作成する場合はcolumn名は大文字小文字も合わせる必要がある。
// 合わせないとselectしたときにバインドできない。
type User struct {
	Id       string `gorm:"primaryKey;column:ID"`
	Password string `gorm:"column:PASSWORD"`
	Name     string `gorm:"column:NAME"`
	SeqNo    int    `gorm:"column:SEQNO"`
}

func (u User) TableName() string {
	return "t_user"
}
