package dto

type RestUser struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SeqNo    int    `json:"seqNo"`
}
