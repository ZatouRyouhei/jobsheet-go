package model

// 既存テーブルに合わせてモデルを作成する場合はcolumn名は大文字小文字も合わせる必要がある。
// 合わせないとselectしたときにバインドできない。
type Attachment struct {
	ID         string `gorm:"primaryKey;column:ID"`
	SeqNo      int    `gorm:"primaryKey;column:SEQNO"`
	FileName   string `gorm:"column:FILENAME"`
	AttachFile []byte `gorm:"column:ATTACHFILE;size:70000"`
}

func (a Attachment) TableName() string {
	return "t_attachment"
}
