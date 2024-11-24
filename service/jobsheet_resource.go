package service

import (
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/database/util"
	"jobsheet-go/dto"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm/clause"
)

func RegistJobSheet(c echo.Context) error {
	var restJobSheet = new(dto.RestJobSheet)
	err := c.Bind(restJobSheet)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}

	// トランザクション開始
	tx := database.Db.Begin()

	newflg := false
	// 新規登録の時はIDを自動採番
	if restJobSheet.ID == "" {
		nextId, err := util.GetNextId(tx)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, "bad request")
		}
		restJobSheet.ID = nextId
		newflg = true
	}

	targetJobSheet := new(model.JobSheet)
	targetJobSheet.ID = restJobSheet.ID
	targetJobSheet.ClientID = restJobSheet.ClientID
	targetJobSheet.BusinessSystemID = restJobSheet.BusinessSystemID
	targetJobSheet.InquiryID = restJobSheet.InquiryID
	targetJobSheet.Department = restJobSheet.Department
	targetJobSheet.Person = restJobSheet.Person
	if restJobSheet.OccurDate != "" && restJobSheet.OccurTime != "" {
		occurDateTime := restJobSheet.OccurDate + " " + restJobSheet.OccurTime
		jst, _ := time.LoadLocation("Asia/Tokyo")
		targetJobSheet.OccurDateTime, err = time.ParseInLocation("2006-01-02 15:04", occurDateTime, jst)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, "bad request")
		}
	}
	targetJobSheet.Title = restJobSheet.Title
	targetJobSheet.Content = restJobSheet.Content
	targetJobSheet.ContactID = restJobSheet.ContactID
	if restJobSheet.LimitDate != "" {
		targetJobSheet.LimitDate, err = time.Parse("2006-01-02", restJobSheet.LimitDate)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, "bad request")
		}
	}
	if restJobSheet.DealID != "" {
		targetJobSheet.DealID = restJobSheet.DealID
	}
	if restJobSheet.CompleteDate != "" {
		targetJobSheet.CompleteDate, err = time.Parse("2006-01-02", restJobSheet.CompleteDate)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, "bad request")
		}
	}
	targetJobSheet.Support = restJobSheet.Support
	targetJobSheet.ResponseTime = restJobSheet.ResponseTime
	if newflg {
		// 新規登録
		result := tx.Create(targetJobSheet)
		if result.Error != nil {
			slog.Error("Error", slog.Any("error", result.Error))
			tx.Rollback()
			return c.String(http.StatusBadRequest, "bad request")
		} else {
			tx.Commit()
		}
	} else {
		// 更新
		result := tx.Save(targetJobSheet)
		if result.Error != nil {
			slog.Error("Error", slog.Any("error", result.Error))
			tx.Rollback()
			return c.String(http.StatusBadRequest, "bad request")
		} else {
			tx.Commit()
		}
	}
	return c.String(http.StatusOK, targetJobSheet.ID)
}

func DeleteJobSheet(c echo.Context) error {
	id := c.Param("id")
	targetJobSheet := new(model.JobSheet)
	targetJobSheet.ID = id
	database.Db.Delete(targetJobSheet)
	// 添付ファイルも削除する。
	database.Db.Where("id = ?", id).Delete(&model.Attachment{})
	return c.String(http.StatusOK, "delete jobsheet")
}

func GetJobSheet(c echo.Context) error {
	id := c.Param("id")
	var targetJobSheet model.JobSheet
	//result := database.Db.Where("id = ?", id).Preload("Client").Preload("Contact").Preload("BusinessSystem").Preload("Deal").Preload("Inquiry").Preload("BusinessSystem.Business").Find(&targetJobSheet)
	// 上記のようにひとつづつPreloadを指定することもできるが、clause.Associationsですべて指定することができる。ただしネストしているものは個別に指定する。
	result := database.Db.Where("id = ?", id).Preload(clause.Associations).Preload("BusinessSystem.Business").Find(&targetJobSheet)
	if result.RowsAffected == 0 {
		// 結果が取得できなかった場合
		return c.String(http.StatusNotFound, "not found")
	} else {
		// 結果が取得できた場合
		// 関連する添付ファイルのファイル名リストを取得する。
		var targetAttachmentList []model.Attachment
		database.Db.Where("id = ?", id).Find(&targetAttachmentList)
		// 返答用の日誌データを作成
		restJobSheet := dto.NewRestSearchJobSheet(targetJobSheet, targetAttachmentList)
		return c.JSON(http.StatusCreated, *restJobSheet)
	}
}

func SearchJobSheet(c echo.Context) error {
	condition := new(dto.RestSearchConditionJobSheet)
	err := c.Bind(condition)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	jobSheetList := util.GetJobSheetList(condition)
	var restSearchJobSheetList []dto.RestSearchJobSheet
	for _, jobSheet := range jobSheetList {
		// 結果が取得できた場合
		// 関連する添付ファイルのファイル名リストを取得する。
		var targetAttachmentList []model.Attachment
		database.Db.Where("id = ?", jobSheet.ID).Find(&targetAttachmentList)
		restSearchJobSheetList = append(restSearchJobSheetList, *dto.NewRestSearchJobSheet(jobSheet, targetAttachmentList))
	}
	return c.JSON(http.StatusCreated, restSearchJobSheetList)
}

func DownloadJobSheet(c echo.Context) error {
	condition := new(dto.RestSearchConditionJobSheet)
	err := c.Bind(condition)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	jobSheetList := util.GetJobSheetList(condition)
	// 自ファイルからの相対バスだとファイルが見つからない。
	// ビルドで生成されるバイナリファイルからの相対パスを指定する。
	f, err := excelize.OpenFile("template/template.xlsx")
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "fileOpenerror")
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Error", slog.Any("error", err))
		}
	}()
	// シート名
	sheetName := "Sheet1"
	// 今日の日付
	today := time.Now()
	// セルスタイル
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
	})
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	for i, jobSheet := range jobSheetList {
		// セルスタイル設定
		f.SetCellStyle(sheetName, "A"+strconv.Itoa(i+3), "Q"+strconv.Itoa(i+3), style)
		// 番号
		f.SetCellValue(sheetName, "A"+strconv.Itoa(i+3), jobSheet.ID)
		// ステータス
		status := ""
		if jobSheet.CompleteDate.IsZero() {
			if jobSheet.LimitDate.Before(today) {
				// 期限超過
				status = "期限超過"
			} else {
				// 期限前
				diffDays := jobSheet.LimitDate.Sub(today).Hours() / 24
				if diffDays <= 3 {
					status = "あと" + strconv.Itoa(int(diffDays)) + "日"
				}
			}
		} else {
			status = "完了"
		}
		f.SetCellValue(sheetName, "B"+strconv.Itoa(i+3), status)
		// 顧客
		f.SetCellValue(sheetName, "C"+strconv.Itoa(i+3), jobSheet.Client.Name)
		// 業務
		f.SetCellValue(sheetName, "D"+strconv.Itoa(i+3), jobSheet.BusinessSystem.Business.Name)
		// システム
		f.SetCellValue(sheetName, "E"+strconv.Itoa(i+3), jobSheet.BusinessSystem.Name)
		// 問合せ区分
		f.SetCellValue(sheetName, "F"+strconv.Itoa(i+3), jobSheet.Inquiry.Name)
		// 部署
		f.SetCellValue(sheetName, "G"+strconv.Itoa(i+3), jobSheet.Department)
		// 担当者
		f.SetCellValue(sheetName, "H"+strconv.Itoa(i+3), jobSheet.Person)
		// 発生日時
		occurDateTime := ""
		if !jobSheet.OccurDateTime.IsZero() {
			occurDateTime = jobSheet.OccurDateTime.Format("2006/01/02 15:04")
		}
		f.SetCellValue(sheetName, "I"+strconv.Itoa(i+3), occurDateTime)
		// 窓口
		f.SetCellValue(sheetName, "J"+strconv.Itoa(i+3), jobSheet.Contact.Name)
		// タイトル
		f.SetCellValue(sheetName, "K"+strconv.Itoa(i+3), jobSheet.Title)
		// 内容
		f.SetCellValue(sheetName, "L"+strconv.Itoa(i+3), jobSheet.Content)
		// 完了期限
		limitDate := ""
		if !jobSheet.LimitDate.IsZero() {
			limitDate = jobSheet.LimitDate.Format("2006/01/02")
		}
		f.SetCellValue(sheetName, "M"+strconv.Itoa(i+3), limitDate)
		// 対応詳細
		f.SetCellValue(sheetName, "N"+strconv.Itoa(i+3), jobSheet.Support)
		// 対応者
		f.SetCellValue(sheetName, "O"+strconv.Itoa(i+3), jobSheet.Deal.Name)
		// 完了日
		completeDate := ""
		if !jobSheet.CompleteDate.IsZero() {
			completeDate = jobSheet.CompleteDate.Format("2006/01/02")
		}
		f.SetCellValue(sheetName, "P"+strconv.Itoa(i+3), completeDate)
		// 対応時間
		f.SetCellValue(sheetName, "Q"+strconv.Itoa(i+3), jobSheet.ResponseTime)
	}

	buf, _ := f.WriteToBuffer()
	response := c.Response()
	response.Writer.Header().Set("Content-Disposition", "attachment; filename=業務日誌.xlsx")
	return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func PdfJobSheet(c echo.Context) error {
	id := c.Param("id")
	var targetJobSheet model.JobSheet
	//result := database.Db.Where("id = ?", id).Preload("Client").Preload("Contact").Preload("BusinessSystem").Preload("Deal").Preload("Inquiry").Preload("BusinessSystem.Business").Find(&targetJobSheet)
	// 上記のようにひとつづつPreloadを指定することもできるが、clause.Associationsですべて指定することができる。ただしネストしているものは個別に指定する。
	result := database.Db.Where("id = ?", id).Preload(clause.Associations).Preload("BusinessSystem.Business").Find(&targetJobSheet)
	if result.RowsAffected == 0 {
		// 結果が取得できなかった場合
		return c.String(http.StatusNotFound, "not found")
	} else {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{
			PageSize: *gopdf.PageSizeA4,
		})
		pdf.AddPage()
		err := pdf.AddTTFFont("genju", "font/GenJyuuGothic-Regular.ttf")
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
		}
		err = pdf.SetFont("genju", "", 10)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
		}
		// 出力日
		rect := gopdf.Rect{W: 80, H: 20}
		pdf.SetX(460)
		pdf.SetY(60)
		op := gopdf.CellOption{
			Align: gopdf.Right | gopdf.Bottom,
		}
		now := time.Now()
		pdf.CellWithOption(&rect, now.Format("2006年01月02日"), op)
		// ID
		drawText(&pdf, 30, 30, targetJobSheet.ID)
		// 顧客
		drawText(&pdf, 60, 110, targetJobSheet.Client.Name)
		// 業務
		drawText(&pdf, 230, 110, targetJobSheet.BusinessSystem.Business.Name)
		// システム
		drawText(&pdf, 390, 110, targetJobSheet.BusinessSystem.Name)
		// 問合せ区分
		drawText(&pdf, 60, 160, targetJobSheet.Inquiry.Name)
		// 部署
		drawText(&pdf, 230, 160, targetJobSheet.Department)
		// 担当者
		drawText(&pdf, 390, 160, targetJobSheet.Person)
		// 発生日
		drawText(&pdf, 60, 220, targetJobSheet.OccurDateTime.Format("2006年01月02日 15時04分"))
		// 窓口
		drawText(&pdf, 230, 220, targetJobSheet.Contact.Name)
		// タイトル
		drawText(&pdf, 60, 270, targetJobSheet.Title)
		// 内容
		// 改行の可能性がある場合は、pdf.MultiCellもしくはpdf.MultiCellWithOptionを使用する。
		// op := gopdf.CellOption{
		// 	Align: gopdf.Left,
		// 	// セルの幅にテキストがおさまらないときの挙動 pdf.MultiCellWithOptionを使用するときのオプション。
		// 	BreakOption: &gopdf.BreakOption{
		// 		// 単語の途中でも改行するモード
		// 		Mode: gopdf.BreakModeStrict,
		// 		// BreakModeStrictの場合で単語の途中で改行される場合のセパレータ文字列
		// 		Separator: "-",
		// 		// 単語の途中では改行しないモード
		// 		// Mode:           gopdf.BreakModeIndicatorSensitive,
		// 		// BreakModeIndicatorSensitiveの場合に単語の区切りとなる文字を指定
		// 		// BreakIndicator: ' ',
		// 	},
		// }
		// テキストの途中に改行コードが入っている場合の処理
		rect = gopdf.Rect{W: 100, H: 20}
		pdf.SetX(60)
		pdf.SetY(330)
		contents := strings.Split(targetJobSheet.Content, "\n")
		for _, content := range contents {
			pdf.MultiCell(&rect, content)
		}
		// 完了期限
		if !targetJobSheet.LimitDate.IsZero() {
			drawText(&pdf, 60, 500, targetJobSheet.LimitDate.Format("2006年01月02日"))
		}
		// 対応詳細
		rect = gopdf.Rect{W: 480, H: 20}
		pdf.SetX(60)
		pdf.SetY(560)
		supports := strings.Split(targetJobSheet.Support, "\n")
		for _, support := range supports {
			pdf.MultiCell(&rect, support)
		}
		// 対応者
		drawText(&pdf, 60, 730, targetJobSheet.Deal.Name)
		// 完了日
		if !targetJobSheet.CompleteDate.IsZero() {
			drawText(&pdf, 230, 730, targetJobSheet.CompleteDate.Format("2006年01月02日"))
		}
		// 対応時間
		drawText(&pdf, 390, 730, strconv.FormatFloat(targetJobSheet.ResponseTime, 'f', -1, 64))

		A4 := *gopdf.PageSizeA4
		A4Tate := gopdf.Rect{W: A4.W, H: A4.H}
		// 引数の2つめはテンプレートファイルのページ番号
		tp1 := pdf.ImportPage("template/jobSheet.pdf", 1, "/MediaBox")
		pdf.UseImportedTemplate(tp1, 0, 0, A4.W, A4.H)

		// 改ページする場合は新たにテンプレートのページを追加する。
		// pdf.AddPage()
		// tp1 = pdf.ImportPage("template/jobSheet.pdf", 1, "/MediaBox")
		// pdf.UseImportedTemplate(tp1, 0, 0, A4.W, A4.H)

		// 位置合わせ用　完成後にコメントアウトする
		drawGrid(&pdf, &A4Tate)

		response := c.Response()
		response.Writer.Header().Set("Content-Disposition", "attachment; filename=業務日誌.pdf")
		return c.Blob(http.StatusOK, "application/pdf", pdf.GetBytesPdf())
	}
}

func drawText(pdf *gopdf.GoPdf, x float64, y float64, s string) {
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, s)
}

func drawGrid(pdf *gopdf.GoPdf, page *gopdf.Rect) {
	ww := 10.0
	for i := 1; i < int(page.W/ww); i++ {
		if i%10 == 0 {
			pdf.SetLineWidth(0.8)
			pdf.SetStrokeColor(50, 50, 100)
		} else {
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(100, 100, 130)
		}
		x := float64(i) * ww
		pdf.Line(x, 0, x, page.H)
	}
	for i := 1; i < int(page.H/ww); i++ {
		if i%10 == 0 {
			pdf.SetLineWidth(0.8)
			pdf.SetStrokeColor(50, 50, 100)
		} else {
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(100, 100, 130)
		}
		y := float64(i) * ww
		pdf.Line(0, y, page.W, y)
	}
}
