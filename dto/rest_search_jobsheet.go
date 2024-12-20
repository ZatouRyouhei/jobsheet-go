package dto

import (
	"jobsheet-go/database/model"
)

type RestSearchJobSheet struct {
	ID             string             `json:"id"`
	Client         RestClient         `json:"client"`
	BusinessSystem RestBusinessSystem `json:"businessSystem"`
	Inquiry        RestInquiry        `json:"inquiry"`
	Department     string             `json:"department"`
	Person         string             `json:"person"`
	OccurDate      string             `json:"occurDate"`
	OccurTime      string             `json:"occurTime"`
	Title          string             `json:"title"`
	Content        string             `json:"content"`
	Contact        RestUser           `json:"contact"`
	LimitDate      string             `json:"limitDate"`
	Deal           RestUser           `json:"deal"`
	CompleteDate   string             `json:"completeDate"`
	Support        string             `json:"support"`
	ResponseTime   float64            `json:"responseTime"`
	FileList       []RestAttachment   `json:"fileList"`
}

func NewRestSearchJobSheet(jobSheet model.JobSheet, attachmentList []model.Attachment) *RestSearchJobSheet {
	restSearchJobSheet := new(RestSearchJobSheet)
	restSearchJobSheet.ID = jobSheet.ID
	restSearchJobSheet.Client = RestClient{
		ID:   jobSheet.Client.ID,
		Name: jobSheet.Client.Name,
	}
	restSearchJobSheet.BusinessSystem = RestBusinessSystem{
		ID:   jobSheet.BusinessSystem.BusinessID,
		Name: jobSheet.BusinessSystem.Name,
		Business: RestBusiness{
			ID:   jobSheet.BusinessSystem.Business.ID,
			Name: jobSheet.BusinessSystem.Business.Name,
		},
	}
	restSearchJobSheet.Inquiry = RestInquiry{
		ID:   jobSheet.Inquiry.ID,
		Name: jobSheet.Inquiry.Name,
	}
	restSearchJobSheet.Department = jobSheet.Department
	restSearchJobSheet.Person = jobSheet.Person
	restSearchJobSheet.OccurDate = jobSheet.OccurDateTime.Format("2006-01-02")
	restSearchJobSheet.OccurTime = jobSheet.OccurDateTime.Format("15:04")
	restSearchJobSheet.Title = jobSheet.Title
	restSearchJobSheet.Content = jobSheet.Content
	restSearchJobSheet.Contact = RestUser{
		Id:       jobSheet.Contact.Id,
		Password: "",
		Name:     jobSheet.Contact.Name,
		SeqNo:    jobSheet.Contact.SeqNo,
	}
	if jobSheet.LimitDate != nil {
		restSearchJobSheet.LimitDate = jobSheet.LimitDate.Format("2006-01-02")
	} else {
		restSearchJobSheet.LimitDate = ""
	}
	restSearchJobSheet.Deal = RestUser{
		Id:       jobSheet.Deal.Id,
		Password: "",
		Name:     jobSheet.Deal.Name,
		SeqNo:    jobSheet.Deal.SeqNo,
	}
	if jobSheet.CompleteDate != nil {
		restSearchJobSheet.CompleteDate = jobSheet.CompleteDate.Format("2006-01-02")
	} else {
		restSearchJobSheet.CompleteDate = ""
	}
	restSearchJobSheet.Support = jobSheet.Support
	restSearchJobSheet.ResponseTime = jobSheet.ResponseTime
	restSearchJobSheet.FileList = []RestAttachment{}
	for _, attachment := range attachmentList {
		restSearchJobSheet.FileList = append(restSearchJobSheet.FileList, RestAttachment{SeqNo: attachment.SeqNo, FileName: attachment.FileName})
	}

	return restSearchJobSheet
}
