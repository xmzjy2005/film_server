package system

import (
	"film_server/config"
	"film_server/plugin/db"
	"gorm.io/gorm"
	"log"
)

type FailureRecord struct {
	gorm.Model
	OriginId    string       `json:"originId"`
	OriginName  string       `json:"originName"`
	Uri         string       `json:"uri"`
	CollectType ResourceType `json:"collectType"`
	PageNumber  int          `json:"pageNumber"`
	Hour        int          `json:"hour"`
	Cause       string       `json:"cause"`  //失败原因
	Status      int          `json:"status"` //重试状态
}

// 设置表明 failure_records
func (fr FailureRecord) TableName() string {
	return config.FailureRecordTableName
}

func SaveFailureRecord(fr FailureRecord) {
	err := db.Mdb.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(fr).Error
		if err != nil {
			log.Println("save failure records data error", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Save failure record affairs failed:", err)
	}
}
