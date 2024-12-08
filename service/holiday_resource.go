package service

import (
	"bufio"
	"fmt"
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func GetHolidayList(c echo.Context) error {
	var holidayList []model.Holiday
	database.Db.Order("HOLIDAY DESC").Find(&holidayList)
	var restHolidayList []dto.RestHoliday
	for _, holiday := range holidayList {
		restHolidayList = append(restHolidayList, dto.NewRestHoliday(holiday))
	}
	return c.JSON(http.StatusCreated, restHolidayList)
}

func RegistHoliday(c echo.Context) error {
	// ファイルは内閣府のホームページに公開されているCSVファイルの想定
	// 文字コードはSJIS
	// ヘッダーあり
	// ダブルクォーテーションなし
	form, err := c.FormFile("file")
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	src, err := form.Open()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "file cant open")
	}
	defer src.Close()

	// SJISの文字コードとして読み込む
	reader := transform.NewReader(src, japanese.ShiftJIS.NewDecoder())
	s := bufio.NewScanner(reader)

	// エラーリスト
	var errorList []dto.RestErrorMessage

	// データ読み込み
	var holidays []model.Holiday
	rowNum := 1 // 行番号
	for s.Scan() {
		// ヘッダー行は無視する
		if rowNum == 1 {
			rowNum += 1
			continue
		}
		fmt.Println(s.Text())
		result := strings.Split(s.Text(), ",")
		if len(result) != 2 {
			errorList = append(errorList, dto.RestErrorMessage{LineNo: rowNum, ErrorMsg: "フォーマットエラー（日付と祝日名称を入力してください。）"})
			rowNum += 1
			continue
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		holiday, err := time.ParseInLocation("2006/1/2", result[0], jst)
		if err != nil {
			errorList = append(errorList, dto.RestErrorMessage{LineNo: rowNum, ErrorMsg: "日付形式のエラー（yyyy/mm/ddとしてください。）"})
			rowNum += 1
			continue
		}
		name := result[1]
		if utf8.RuneCountInString(name) > 20 {
			errorList = append(errorList, dto.RestErrorMessage{LineNo: rowNum, ErrorMsg: "祝日名称エラー（祝日名称は20文字以内としてください。）"})
			rowNum += 1
			continue
		}
		holidays = append(holidays, model.Holiday{
			Holiday: holiday,
			Name:    name,
		})
		rowNum += 1
	}

	// 入力データ0件の時はエラー
	if rowNum-1 < 2 {
		errorList = append(errorList, dto.RestErrorMessage{LineNo: 0, ErrorMsg: "ヘッダーを含め2行以上入力してください。"})
	}

	// エラーがないときは登録処理を実施
	// 日付が重複している行は更新する。
	if len(errorList) == 0 {
		database.Db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "HOLIDAY"}},
			DoUpdates: clause.AssignmentColumns([]string{"NAME"}),
		}).Create(&holidays)
	}
	return c.JSON(http.StatusCreated, errorList)
}
