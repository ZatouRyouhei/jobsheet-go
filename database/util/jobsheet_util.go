package util

import (
	"fmt"
	"jobsheet-go/database"
	"jobsheet-go/database/model"
	"jobsheet-go/dto"
	"regexp"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetNextId(tx *gorm.DB) (string, error) {
	now := time.Now()
	idHeader := now.Format("2006-01")
	var maxJobSheet model.JobSheet
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id LIKE ?", idHeader+"%").Order("id desc").First(&maxJobSheet)
	nextId := "001"
	if result.RowsAffected > 0 {
		maxSeqNo, err := strconv.Atoi(maxJobSheet.ID[8:])
		if err != nil {
			return "", err
		}
		nextId = fmt.Sprintf("%03d", maxSeqNo+1)
	}
	return idHeader + "-" + nextId, nil
}

func GetJobSheetList(condition *dto.RestSearchConditionJobSheet) []model.JobSheet {
	var jobSheetList []model.JobSheet
	stat := database.Db
	if condition.Client != 0 {
		stat = stat.Where("CLIENT_ID = ?", condition.Client)
	}
	if condition.Business != 0 {
		//stat.InnerJoins("t_business", stat.Where(&model.BusinessSystem{BusinessID: condition.Business}))
		stat = stat.Joins("join t_business_system on t_job_sheet.BUSINESSSYSTEM_ID = t_business_system.ID").Where("t_business_system.BUSINESS_ID = ?", condition.Business)
	}
	if condition.BusinessSystem != 0 {
		stat = stat.Where("BUSINESSSYSTEM_ID = ?", condition.BusinessSystem)
	}
	if condition.Inquiry != 0 {
		stat = stat.Where("INQUIRY_ID = ?", condition.Inquiry)
	}
	if condition.Contact != "" {
		stat = stat.Where("CONTACT_ID = ?", condition.Contact)
	}
	if condition.Deal != "" {
		stat = stat.Where("DEAL_ID = ?", condition.Deal)
	}
	if condition.OccurDateFrom != "" {
		stat = stat.Where("OCCURDATETIME >= STR_TO_DATE(?, '%Y-%m-%d')", condition.OccurDateFrom)
	}
	if condition.OccurDateTo != "" {
		stat = stat.Where("OCCURDATETIME < DATE_ADD(STR_TO_DATE(?, '%Y-%m-%d'), INTERVAL 1 DAY)", condition.OccurDateTo)
	}
	if condition.CompleteSign == 1 {
		stat = stat.Where("COMPLETEDATE IS NOT NULL")
	}
	if condition.CompleteSign == 2 {
		stat = stat.Where("COMPLETEDATE IS NULL")
	}
	if condition.LimitDate != "" {
		stat = stat.Where("LIMITDATE <= STR_TO_DATE(?, '%Y-%m-%d')", condition.LimitDate)
	}
	if condition.Keyword != "" {
		reg := "( |　)+"
		keywordArr := regexp.MustCompile(reg).Split(condition.Keyword, -1)
		keywordCond := database.Db.Table("t_job_sheet")
		for _, keyword := range keywordArr {
			keywordCond = keywordCond.Or("TITLE LIKE ?", "%"+keyword+"%")
			keywordCond = keywordCond.Or("CONTENT LIKE ?", "%"+keyword+"%")
			keywordCond = keywordCond.Or("SUPPORT LIKE ?", "%"+keyword+"%")
			keywordCond = keywordCond.Or("DEPARTMENT LIKE ?", "%"+keyword+"%")
			keywordCond = keywordCond.Or("PERSON LIKE ?", "%"+keyword+"%")
		}
		stat = stat.Where(keywordCond)
	}
	stat.Preload(clause.Associations).Preload("BusinessSystem.Business").Order("ID desc").Find(&jobSheetList)
	return jobSheetList
}
