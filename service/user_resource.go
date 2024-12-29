package service

import (
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	user := new(dto.RestUser)
	err := c.Bind(user)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	var db_user model.User
	result := database.Db.Where("id = ? and password = ?", user.Id, user.Password).First(&db_user)
	if result.RowsAffected == 1 {
		// JWT生成
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		})
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
		}
		return c.JSON(http.StatusOK, dto.RestLoginUser{
			User: dto.RestUser{
				Id:       db_user.Id,
				Password: "",
				Name:     db_user.Name,
				SeqNo:    db_user.SeqNo,
			},
			Token: t,
		})
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
	database.Db.Order("seqno").Find(&dbUserList)
	var rUserList []dto.RestUser
	for _, user := range dbUserList {
		rUser := dto.NewRestUser(user)
		rUserList = append(rUserList, rUser)
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
	var dbUser model.User
	result := database.Db.Where("id = ?", user.Id).Find(&dbUser)
	if result.RowsAffected == 0 {
		// 新規登録
		dbUser.Id = user.Id
		dbUser.Name = user.Name
		dbUser.Password = user.Password
		// 連番取得
		var maxSeqUser model.User
		database.Db.Order("seqno desc").First(&maxSeqUser)
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
	var jobSheetData model.JobSheet
	result := database.Db.Where("contact_id = ? or deal_id = ?", id, id).Find(&jobSheetData)
	if result.RowsAffected > 0 {
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
		targetId := userList[seqNo-1]
		var targetUser model.User
		database.Db.Where("id = ?", targetId).First(&targetUser)
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
	result := database.Db.Where("id = ?", restUser.Id).First(&targetUser)
	if result.RowsAffected == 1 {
		targetUser.Password = restUser.Password
		database.Db.Save(&targetUser)
		return nil
	} else {
		return c.String(http.StatusBadRequest, "bad request")
	}
}
