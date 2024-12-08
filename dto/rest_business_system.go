package dto

import "jobsheet-go/database/model"

type RestBusinessSystem struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Business RestBusiness `json:"business"`
}

func NewRestBusinessSystem(businessSystem model.BusinessSystem) RestBusinessSystem {
	restBusinessSystem := new(RestBusinessSystem)
	restBusinessSystem.ID = businessSystem.ID
	restBusinessSystem.Name = businessSystem.Name
	restBusinessSystem.Business = NewRestBusiness(businessSystem.Business)
	return *restBusinessSystem
}

func NewBusinessSystem(restBusinessSystem RestBusinessSystem) model.BusinessSystem {
	businessSystem := new(model.BusinessSystem)
	businessSystem.ID = restBusinessSystem.ID
	businessSystem.Name = restBusinessSystem.Name
	businessSystem.BusinessID = restBusinessSystem.Business.ID
	businessSystem.Business = NewBusiness(restBusinessSystem.Business)
	return *businessSystem
}
