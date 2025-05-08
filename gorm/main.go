package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMNTCaseInfo_Equal2() {
	db, err := gorm.Open(
		mysql.Open("root:my-secret-pw@tcp(xxx.xxx.xx.xx:13306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{},
	)
	if err != nil {
		log.Fatal(err)
	}

	db = db.Begin()
	defer db.Commit()

	tx := db
	isTX := false // 入参是否已经是事务，避免重复开启事务
	if committer, ok := db.Statement.ConnPool.(gorm.TxCommitter); ok && committer != nil {
		isTX = true
		tx.SavePoint("UpdateGenerationModule")
	} else {
		tx = db.Begin()
	}
	defer func() {
		if err != nil {
			if isTX {
				tx.RollbackTo("UpdateGenerationModule")
				return
			}
			tx.Rollback()
			return
		}
		if isTX {
			return
		}
		if e := tx.Commit().Error; e != nil {
			err = e
		}
	}()

	if err := tx.Error; err != nil {
		log.Fatal(err)
	}
	type User struct {
		ID   int64  `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	err = tx.Raw("select * from user where id=123456789").Scan(&User{}).Error
	if err != nil {
		log.Fatalf("err:%v", err)
	}
}
