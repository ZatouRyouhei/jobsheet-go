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

func GetClientList(c echo.Context) error {
	var clientList []model.Client
	database.Db.Order("ID").Find(&clientList)
	var restClientList []dto.RestClient
	for _, client := range clientList {
		restClientList = append(restClientList, dto.NewClient(client))
	}
	return c.JSON(http.StatusCreated, restClientList)
}

func RegistClient(c echo.Context) error {
	restClient := new(dto.RestClient)
	err := c.Bind(restClient)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	targetClient := new(model.Client)
	targetClient.ID = restClient.ID
	targetClient.Name = restClient.Name
	database.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ID"}},
		DoUpdates: clause.AssignmentColumns([]string{"NAME"}),
	}).Create(targetClient)

	return c.String(http.StatusOK, "regist client")
}

func DeleteClient(c echo.Context) error {
	id := c.Param("id")
	resultFlg := "0"
	// 業務日誌で使用されているか確認
	var checkJobSheetList []model.JobSheet
	checkResult := database.Db.Where("CLIENT_ID = ?", id).Find(&checkJobSheetList)
	if checkResult.RowsAffected > 0 {
		// 使用されている場合は削除しない
		resultFlg = "1"
	} else {
		// 使用されていない場合は削除する。
		var targetClient model.Client
		database.Db.Where("ID = ?", id).Delete(&targetClient)
	}
	return c.String(http.StatusOK, resultFlg)
}
