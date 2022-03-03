package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId       string //`json:"user_id`
	UserName     string //`json:"user_name`
	UserPassword string //`json:"user_password`
}

type Meeting struct {
	MeetingId        int    `gorm:"AUTO_INCREMENT"`
	MeetingName      string //`json:"meeting_name`
	MeetingStartTime time.Time
}

type Participant struct {
	MeetingId        int    //`json:"meeting_id`
	UserId           string //`json:"user_id`
	SpeakNum         int
	ParticipantOrder int
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

func signupUser(db *gorm.DB, userId string, userName string, userPassword string) bool {
	user := User{UserId: userId, UserName: userName, UserPassword: userPassword}
	if err := db.Create(&user).Error; err == nil {
		fmt.Printf("signup成功: %s, %s, %s", userId, userName, userPassword)
		return true
	} else {
		fmt.Println("signup失敗")
		return false
	}
}

func loginUser(db *gorm.DB, userId string, userPassword string) bool {
	var user User
	// err := db.Find(&user).Error
	err := db.First(&user, "user_id = ? AND user_password = ?", userId, userPassword).Error
	if err == nil {
		fmt.Printf("login成功: %s, %s\n", userId, userPassword)
		return true
	} else {
		fmt.Println("login失敗")
		return false
	}
}

func createMeeting(db *gorm.DB, meetingName string, startTimeStr string, presenters []string) (int, string, string, []string) {
	var (
		user         User
		layout       = "2006/01/02 15:04:05"
		startTime, _ = time.Parse(layout, startTimeStr)
		meeting      = Meeting{MeetingName: meetingName, MeetingStartTime: startTime}
	)

	if err := db.Create(&meeting).Error; err == nil {
		for i, presenter := range presenters {
			if err := db.First(&user, "user_id = ?", presenter).Error; err == nil {
				participant := Participant{MeetingId: meeting.MeetingId, UserId: user.UserId, SpeakNum: 0, ParticipantOrder: i}
				if err := db.Create(&participant).Error; err != nil { // TODO: transaction
					fmt.Printf("create失敗(発表者%sの登録に失敗しました): %s, %s, %s\n", presenter, meetingName, startTimeStr, presenters)
					return -1, "", "", []string{}
				}
			} else {
				fmt.Printf("create失敗(発表者%sが見つかりません): %s, %s, %s\n", presenter, meetingName, startTimeStr, presenters)
				return -1, "", "", []string{}
			}
		}
		fmt.Printf("create成功: %s, %s, %s\n", meetingName, startTimeStr, presenters)
		return meeting.MeetingId, meetingName, startTimeStr, presenters
	} else {
		fmt.Printf("create失敗(会議の登録に失敗しました): %s, %s, %s\n", meetingName, startTimeStr, presenters)
		return -1, "", "", []string{}
	}
}
