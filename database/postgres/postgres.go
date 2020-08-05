/*
 * @Version: 0.0.1
 * @Author: ider
 * @Date: 2020-06-02 17:11:59
 * @LastEditors: ider
 * @LastEditTime: 2020-08-05 14:14:10
 * @Description:
 */
package postgres

import (
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var (
	WikiDB *gorm.DB
)

type WikiWordCount struct {
	ID             uint   `gorm:"primaryKey"`
	Title          string `gorm:"uniqueIndex"`
	CharCount      int32
	WordCount      int32
	CleanWordCount int32
}

func NewWikiDBConn(PostgresConn string) {
	var err error
	WikiDB, err = gorm.Open(postgres.Open(PostgresConn), &gorm.Config{})
	if err != nil {
		log.Fatal("ping 失败", err)
	}
	WikiDB.AutoMigrate(&WikiWordCount{})

}
