package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId       string //`json:"user_id`
	UserName     string //`json:"user_name`
	UserPassword string //`json:"user_password`
}

// SQLConnect DB接続
func sqlConnect() (database *gorm.DB, err error) {
	DBMS := os.Getenv("DBMS")
	USER := os.Getenv("USER")
	PASS := os.Getenv("PASS")
	PROTOCOL := os.Getenv("PROTOCOL")
	DBNAME := os.Getenv("DBNAME")

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}

func connectDB() *gorm.DB {
	// DB接続
	db, err := sqlConnect()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("DBへの接続に成功しました")
	}

	return db
}

func findUser(db *gorm.DB) User {
	var user User
	err := db.Find(&user).Error
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(db.HasTable("users"))
	fmt.Println(user.UserId, user.UserName, user.UserPassword)

	return user
}

func addUser(db *gorm.DB) []User {
	var testuser = User{UserId: "test02", UserName: "test_02", UserPassword: "testpass02"}
	err0 := db.Create(&testuser).Error
	if err0 != nil {
		panic(err0.Error())
	}
	var users []User
	err := db.Find(&users).Error
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(db.HasTable("users"))
	for _, user := range users {
		fmt.Println(user.UserId, user.UserName, user.UserPassword)
	}

	return users
}
