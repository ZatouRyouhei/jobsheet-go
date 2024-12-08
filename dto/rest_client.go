package dto

import "jobsheet-go/database/model"

type RestClient struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewClient(client model.Client) RestClient {
	restClient := new(RestClient)
	restClient.ID = client.ID
	restClient.Name = client.Name
	return *restClient
}
