package dto

import "jobsheet-go/database/model"

type RestUser struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SeqNo    int    `json:"seqNo"`
}

func NewRestUser(user model.User) RestUser {
	restUser := new(RestUser)
	restUser.Id = user.Id
	restUser.Password = user.Password
	restUser.Name = user.Name
	restUser.SeqNo = user.SeqNo
	return *restUser
}
