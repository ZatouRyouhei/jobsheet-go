package constant

const BASE_URL string = "/jobsheet/webresources"

// const DSN string = "postgres://postgres:ryouhei_postgre@localhost:5432/rust_sample?sslmode=disable"
const DSN string = "jobsheet:jobsheet@tcp(localhost:3306)/jobsheetdb?charset=utf8mb4&parseTime=True&loc=Local"

type AttachmentMode int

const (
	DBMode AttachmentMode = iota
	FileMode
)
const ATTACHMENT_MODE AttachmentMode = FileMode

const ATTACHMENT_BASE_DIR string = "jobsheet_attachement/"
