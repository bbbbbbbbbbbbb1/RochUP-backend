package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId       string //`json:"user_id"`
	UserName     string //`json:"user_name"`
	UserPassword string //`json:"user_password"`
}

type Meeting struct {
	MeetingId        int       //`json:"meeting_id"`
	MeetingName      string    //`json:"meeting_name"`
	MeetingStartTime time.Time //`json:"meetingstarttime"`
}

type Participant struct {
	MeetingId        int    //`json:"meeting_id"`
	UserId           string //`json:"user_id"`
	SpeakNum         int    //`json:"speaknum"`
	ParticipantOrder int    //`json:"participantorder"`
}

type ByParticipantOrder []Participant

func (p ByParticipantOrder) Len() int           { return len(p) }
func (p ByParticipantOrder) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByParticipantOrder) Less(i, j int) bool { return p[i].ParticipantOrder < p[j].ParticipantOrder }

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
		fmt.Printf("signup成功: %s, %s, %s\n", userId, userName, userPassword)
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

func joinMeeting(db *gorm.DB, userId string, meetingId int) (bool, string, time.Time, []string) {
	var user User
	var meeting Meeting
	var participant Participant
	participants := make([]Participant, 0, 10)
	user_info := db.First(&user, "user_id = ?", userId)
	meeting_info := db.First(&meeting, "meeting_id = ?", meetingId)
	if user_info.Error == nil && meeting_info.Error == nil {
		participant_info := db.First(&participant, "user_id = ? AND meeting_id = ?", userId, meetingId)
		if participant_info.Error != nil {
			participant.MeetingId = meetingId
			participant.UserId = userId
			participant.SpeakNum = 0
			participant.ParticipantOrder = -1
			if err := db.Create(&participant).Error; err == nil {
				fmt.Printf("参加者追加成功: %s, %d\n", userId, meetingId)
			} else {
				fmt.Println("発表者追加失敗")
				temp_string := []string{"false"}
				return false, "false", time.Now(), temp_string
			}
		}
		participants_err := db.Find(&participants, "meeting_id = ? AND participant_order != -1", meetingId).Error
		if participants_err != nil {
			fmt.Println("join失敗")
			temp_string := []string{"false"}
			return false, "false", time.Now(), temp_string
		}
		presenter_names := make([]string, 0, 10)

		sort.Sort(ByParticipantOrder(participants))

		for _, p := range participants {
			if p.ParticipantOrder != -1 {
				presenter_id := p.UserId
				user_err := db.First(&user, "user_id = ?", presenter_id).Error
				if user_err != nil {
					fmt.Println("join失敗")
					temp_string := []string{"false"}
					return false, "false", time.Now(), temp_string
				}
				presenter_names = append(presenter_names, user.UserName)
			}
		}

		fmt.Printf("join成功: %s, %d\n", userId, meetingId)
		return true, meeting.MeetingName, meeting.MeetingStartTime, presenter_names

	} else {
		fmt.Println("join失敗")
		temp_string := []string{"false"}
		return false, "false", time.Now(), temp_string
	}
}
