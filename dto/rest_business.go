package dto

import (
	"jobsheet-go/database/model"
)

type RestBusiness struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRestBusiness(business model.Business) RestBusiness {
	return RestBusiness{
		ID:   business.ID,
		Name: business.Name,
	}
}

func NewBusiness(restBusiness RestBusiness) model.Business {
	return model.Business{
		ID:   restBusiness.ID,
		Name: restBusiness.Name,
	}
}
