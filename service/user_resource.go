package service

import (
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	user := new(dto.RestUser)
	err := c.Bind(user)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	var cnt int64
	var db_user model.User
	database.Db.Select("id, name, password, seqno").Where("id = ? and password = ?", user.Id, user.Password).First(&db_user).Count(&cnt)
	if cnt == 1 {
		user.Name = db_user.Name
		user.SeqNo = db_user.SeqNo
		return c.JSON(http.StatusCreated, user)
	} else {
		return nil
	}

	// id := c.Param("id")
	// var cnt int64
	// var db_user model.User
	// database.Db.Debug().Select("id, name, seqno").Where("id = ?", id).First(&db_user).Count(&cnt)
	// // fmt.Println("cnt = " + string(cnt))
	// // if cnt == 1 {
	// // 	fmt.Println("id:" + db_user.Id + ",name:" + db_user.Name)
	// // 	return c.String(http.StatusOK, "OK!")
	// // } else {
	// // 	return c.String(http.StatusBadRequest, "record not found")
	// // }
	// fmt.Println(cnt)
	// fmt.Println(db_user)
	// fmt.Println("id:" + db_user.Id + ",name:" + db_user.Name)

	// return c.String(http.StatusOK, "OK!")
}

func GetList(c echo.Context) error {
	var dbUserList []model.User
	database.Db.Select("id, password, name, seqno").Order("seqno").Find(&dbUserList)
	var rUserList []dto.RestUser
	for _, user := range dbUserList {
		rUser := new(dto.RestUser)
		rUser.Id = user.Id
		rUser.Password = user.Password
		rUser.Name = user.Name
		rUser.SeqNo = user.SeqNo
		rUserList = append(rUserList, *rUser)
	}
	return c.JSON(http.StatusCreated, rUserList)
}

func RegistUser(c echo.Context) error {
	user := new(dto.RestUser)
	err := c.Bind(user)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	var cnt int64
	var dbUser model.User
	database.Db.Select("id, password, name, seqno").Where("id = ?", user.Id).Find(&dbUser).Count(&cnt)
	if cnt == 0 {
		// 新規登録
		dbUser.Id = user.Id
		dbUser.Name = user.Name
		dbUser.Password = user.Password
		// 連番取得
		var maxSeqUser model.User
		database.Db.Select("id, password, name, seqno").Order("seqno desc").First(&maxSeqUser)
		nextSeqNo := maxSeqUser.SeqNo + 1
		dbUser.SeqNo = nextSeqNo
		// データベースに登録
		database.Db.Create(&dbUser)
	} else {
		// 更新
		if user.Password != "" {
			dbUser.Password = user.Password
		}
		dbUser.Name = user.Name
		database.Db.Save(&dbUser)
	}
	return c.String(http.StatusOK, "user updated")
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	var cnt int64
	var jobSheetData model.JobSheet
	database.Db.Select("id").Where("contact_id = ? or deal_id = ?", id, id).Find(&jobSheetData).Count(&cnt)
	if cnt > 0 {
		// 業務日誌で使用中のユーザは削除しない。
		return c.String(http.StatusOK, "1")
	} else {
		// 使用されてない場合は削除する。
		targetUser := model.User{
			Id: id,
		}
		database.Db.Delete(&targetUser)
		return c.String(http.StatusOK, "0")
	}
}

func ChangeSeq(c echo.Context) error {
	var restUserList dto.RestUserList
	err := c.Bind(&restUserList)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	userList := restUserList.UserList
	for seqNo := 1; seqNo <= len(userList); seqNo++ {
		user := userList[seqNo-1]
		var targetUser model.User
		database.Db.Select("id, password, name, seqno").Where("id = ?", user.Id).First(&targetUser)
		targetUser.SeqNo = seqNo
		database.Db.Save(&targetUser)
	}
	return c.String(http.StatusOK, "change seq")
}

func ChangePassword(c echo.Context) error {
	var restUser dto.RestUser
	err := c.Bind(&restUser)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	var targetUser model.User
	var cnt int64
	database.Db.Select("id, password, name, seqno").Where("id = ?", restUser.Id).First(&targetUser).Count(&cnt)
	if cnt == 1 {
		targetUser.Password = restUser.Password
		database.Db.Save(&targetUser)
		return nil
	} else {
		return c.String(http.StatusBadRequest, "bad request")
	}
}
