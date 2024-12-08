package service

import (
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetInquiryList(c echo.Context) error {
	var inquiryList []model.Inquiry
	database.Db.Order("ID").Find(&inquiryList)
	var restInquiryList []dto.RestInquiry
	for _, inquiry := range inquiryList {
		restInquiryList = append(restInquiryList, dto.NewRestInquiry(inquiry))
	}
	return c.JSON(http.StatusCreated, restInquiryList)
}
