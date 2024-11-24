package dto

type RestJobSheet struct {
	ID               string  `json:"id"`
	ClientID         int     `json:"client"`
	BusinessID       int     `json:"business"`
	BusinessSystemID int     `json:"businessSystem"`
	InquiryID        int     `json:"inquiry"`
	Department       string  `json:"department"`
	Person           string  `json:"person"`
	OccurDate        string  `json:"occurDate"`
	OccurTime        string  `json:"occurTime"`
	Title            string  `json:"title"`
	Content          string  `json:"content"`
	ContactID        string  `json:"contact"`
	LimitDate        string  `json:"limitDate"`
	DealID           string  `json:"deal"`
	CompleteDate     string  `json:"completeDate"`
	Support          string  `json:"support"`
	ResponseTime     float64 `json:"responseTime"`
}
