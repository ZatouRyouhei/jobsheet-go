package constant

const BASE_URL string = "/jobsheet/webresources"

// const DSN string = "postgres://postgres:ryouhei_postgre@localhost:5432/rust_sample?sslmode=disable"
const DSN string = "jobsheet:jobsheet@tcp(localhost:3306)/jobsheetdb?charset=utf8mb4&parseTime=True&loc=Local"

// 1:DBモード 2:FILEモード
const ATTACHMENT_MODE int = 2

const ATTACHMENT_BASE_DIR string = "C:/jobsheet_attachement/"
