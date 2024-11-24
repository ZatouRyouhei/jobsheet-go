package main

import (
	"jobsheet-go/constant"
	"jobsheet-go/service"

	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	e.POST(constant.BASE_URL+"/user/login", service.GetUser)

	e.GET(constant.BASE_URL+"/user/getList", service.GetList)

	e.POST(constant.BASE_URL+"/user/regist", service.RegistUser)

	e.DELETE(constant.BASE_URL+"/user/delete/:id", service.DeleteUser)

	e.POST(constant.BASE_URL+"/user/changeSeq", service.ChangeSeq)

	e.POST(constant.BASE_URL+"/jobsheet/regist", service.RegistJobSheet)

	e.DELETE(constant.BASE_URL+"/jobsheet/delete/:id", service.DeleteJobSheet)

	e.GET(constant.BASE_URL+"/jobsheet/get/:id", service.GetJobSheet)

	e.POST(constant.BASE_URL+"/jobsheet/search", service.SearchJobSheet)

	e.POST(constant.BASE_URL+"/jobsheet/download", service.DownloadJobSheet)

	e.GET(constant.BASE_URL+"/jobsheet/pdf/:id", service.PdfJobSheet)
}
