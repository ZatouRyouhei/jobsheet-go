package service

import (
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func GetBusinessList(c echo.Context) error {
	var businessList []model.Business
	database.Db.Order("ID").Find(&businessList)
	var restBusinessList []dto.RestBusiness
	for _, business := range businessList {
		restBusinessList = append(restBusinessList, dto.NewRestBusiness(business))
	}
	return c.JSON(http.StatusCreated, restBusinessList)
}

func RegistBusiness(c echo.Context) error {
	restBusiness := new(dto.RestBusiness)
	err := c.Bind(restBusiness)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	targetBusiness := new(model.Business)
	targetBusiness.ID = restBusiness.ID
	targetBusiness.Name = restBusiness.Name
	database.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ID"}},
		DoUpdates: clause.AssignmentColumns([]string{"NAME"}),
	}).Create(targetBusiness)
	return c.String(http.StatusOK, "regist business")
}

func DeleteBusiness(c echo.Context) error {
	id := c.Param("id")
	// システムに紐づいているか確認
	var businessSystem []model.BusinessSystem
	result := database.Db.Where("BUSINESS_ID = ?", id).Find(&businessSystem)
	var resultFlg string
	if result.RowsAffected > 0 {
		// 紐づいているシステムがある場合は削除しない。
		resultFlg = "1"
	} else {
		// 紐づいているシステムがない場合は削除する。
		var business model.Business
		database.Db.Where("ID = ?", id).Delete(&business)
		resultFlg = "0"
	}
	return c.String(http.StatusOK, resultFlg)
}
