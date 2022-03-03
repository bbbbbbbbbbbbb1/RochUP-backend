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
	MeetingId        int    `gorm:"AUTO_INCREMENT"`
	MeetingName      string //`json:"meeting_name`
	MeetingStartTime time.Time
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
			fmt.Println("会議非存在")
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
					fmt.Println("ユーザー非存在")
					temp_string := []string{"false"}
					return false, "false", time.Now(), temp_string
				}
				presenter_names = append(presenter_names, user.UserName)
			}
		}

		fmt.Printf("join成功: %s, %d\n", userId, meetingId)
		return true, meeting.MeetingName, meeting.MeetingStartTime, presenter_names

	} else {
		fmt.Println("ユーザーもしくは会議が非存在")
		temp_string := []string{"false"}
		return false, "false", time.Now(), temp_string
	}
}
