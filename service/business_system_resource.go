package service

import (
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	"log/slog"
	"net/http"

	"gorm.io/gorm/clause"

	"github.com/labstack/echo/v4"
)

func GetBusinessSystemList(c echo.Context) error {
	id := c.Param("id")
	var businessSystemList []model.BusinessSystem
	if id != "" {
		database.Db.Where("BUSINESS_ID = ?", id).Order("ID").Preload("Business").Find(&businessSystemList)
	} else {
		database.Db.Order("ID").Preload("Business").Find(&businessSystemList)
	}
	var restBusinessSystemList []dto.RestBusinessSystem
	for _, businessSystem := range businessSystemList {
		restBusinessSystemList = append(restBusinessSystemList, dto.NewRestBusinessSystem(businessSystem))
	}
	return c.JSON(http.StatusCreated, restBusinessSystemList)
}

func RegistSystem(c echo.Context) error {
	restSystem := new(dto.RestBusinessSystem)
	err := c.Bind(restSystem)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	targetBusinessSystem := new(model.BusinessSystem)
	targetBusinessSystem.ID = restSystem.ID
	targetBusinessSystem.Name = restSystem.Name
	targetBusinessSystem.BusinessID = restSystem.Business.ID
	database.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ID"}},
		DoUpdates: clause.AssignmentColumns([]string{"NAME", "BUSINESS_ID"}),
	}).Create(targetBusinessSystem)

	return c.String(http.StatusOK, "regist business")
}

func DeleteSystem(c echo.Context) error {
	id := c.Param("id")
	resultFlg := "0"
	// 業務日誌で使用されているか確認
	var checkJobSheet []model.JobSheet
	checkResult := database.Db.Where("BUSINESSSYSTEM_ID = ?", id).Find(&checkJobSheet)
	if checkResult.RowsAffected > 0 {
		// 使用されている場合は削除しない
		resultFlg = "1"
	} else {
		// 使用されていない場合は削除する。
		var targetSystem model.BusinessSystem
		database.Db.Where("id = ?", id).Delete(&targetSystem)
	}
	return c.String(http.StatusOK, resultFlg)
}
