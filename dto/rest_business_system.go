package dto

type RestBusinessSystem struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Business RestBusiness `json:"business"`
}
