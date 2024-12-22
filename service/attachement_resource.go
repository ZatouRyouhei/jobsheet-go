package service

import (
	"bytes"
	"io"
	"jobsheet-go/constant"
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	fileUtil "jobsheet-go/util/file"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegistAttachment(c echo.Context) error {
	id := c.Param("id")
	form, err := c.MultipartForm()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	files := form.File["file"]
	for _, file := range files {
		// 連番の最大値を取得する
		maxAttachment := new(model.Attachment)
		result := database.Db.Where("ID = ?", id).Order("SEQNO desc").First(maxAttachment)
		nextSeqNo := 1
		if result.RowsAffected > 0 {
			nextSeqNo = maxAttachment.SeqNo + 1
		}
		if constant.ATTACHMENT_MODE == constant.DBMode {
			// DBモードのとき
			targetAttachement := new(model.Attachment)
			targetAttachement.ID = id
			targetAttachement.SeqNo = nextSeqNo
			targetAttachement.FileName = file.Filename
			src, err := file.Open()
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, "file cant open")
			}
			defer src.Close()
			var buf bytes.Buffer
			io.Copy(&buf, src)
			targetAttachement.AttachFile = buf.Bytes()
			database.Db.Create(targetAttachement)
		} else {
			// FILEモードの時
			// DBにデータを登録する。（添付ファイルはDBには入れずにサーバに保存する。）
			targetAttachement := new(model.Attachment)
			targetAttachement.ID = id
			targetAttachement.SeqNo = nextSeqNo
			targetAttachement.FileName = file.Filename
			database.Db.Create(targetAttachement)

			src, err := file.Open()
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, "file cant open")
			}
			defer src.Close()

			// フォルダに配置する。
			putDir := constant.ATTACHMENT_BASE_DIR + id + "/" + strconv.Itoa(nextSeqNo)
			// フォルダが存在しない場合は作成
			if !fileUtil.IsDir(putDir) {
				err := os.MkdirAll(putDir, os.ModePerm)
				if err != nil {
					slog.Error("Error", slog.Any("error", err))
					return c.String(http.StatusBadRequest, "putDir cant create")
				}
			}
			dst, err := os.Create(putDir + "/" + file.Filename)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, "file cant put")
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, "file cant copy")
			}
		}
	}
	// 返信用の添付ファイルリストを作成
	var attachmentList []model.Attachment
	database.Db.Where("ID = ?", id).Order("SEQNO").Find(&attachmentList)
	var restAttachmentList []dto.RestAttachment
	for _, attachment := range attachmentList {
		restAttachmentList = append(restAttachmentList, dto.RestAttachment{
			SeqNo:    attachment.SeqNo,
			FileName: attachment.FileName,
		})
	}
	return c.JSON(http.StatusCreated, restAttachmentList)
}

func DownloadAttachment(c echo.Context) error {
	id := c.Param("id")
	seqNo := c.Param("seqNo")
	targetAttachment := new(model.Attachment)
	result := database.Db.Where("ID = ? AND SEQNO = ?", id, seqNo).Find(targetAttachment)
	if result.RowsAffected == 0 {
		return c.String(http.StatusBadRequest, "file not exists")
	}

	var targetFile []byte
	if constant.ATTACHMENT_MODE == constant.DBMode {
		// DBモードの時
		targetFile = targetAttachment.AttachFile
	} else {
		// FILEモードの時
		fileLocation := constant.ATTACHMENT_BASE_DIR + id + "/" + seqNo + "/" + targetAttachment.FileName
		file, err := os.ReadFile(fileLocation)
		if err != nil {
			return c.String(http.StatusBadRequest, "file not exists")
		}
		targetFile = file
	}
	response := c.Response()
	response.Writer.Header().Set("Content-Disposition", "attachment; filename="+targetAttachment.FileName)
	return c.Blob(http.StatusOK, "application/octet-stream", targetFile)
}

func DeleteAttachment(c echo.Context) error {
	id := c.Param("id")
	seqNo := c.Param("seqNo")
	database.Db.Where("ID = ? AND SEQNO = ?", id, seqNo).Delete(&model.Attachment{})
	if constant.ATTACHMENT_MODE == constant.FileMode {
		// FILEモードの場合はファイルも削除する
		targetDir := constant.ATTACHMENT_BASE_DIR + id + "/" + seqNo
		err := os.RemoveAll(targetDir)
		if err != nil {
			return c.String(http.StatusBadRequest, "delete failed")
		}
	}

	// 返信用の添付ファイルリストを作成
	var attachmentList []model.Attachment
	database.Db.Where("ID = ?", id).Order("SEQNO").Find(&attachmentList)
	restAttachmentList := []dto.RestAttachment{}
	for _, attachment := range attachmentList {
		restAttachmentList = append(restAttachmentList, dto.RestAttachment{
			SeqNo:    attachment.SeqNo,
			FileName: attachment.FileName,
		})
	}
	return c.JSON(http.StatusCreated, restAttachmentList)
}
