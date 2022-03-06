package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	presenEndMessage       = "発表ありがとうございました．"
	questionBodyAskMessage = "匿名質問です．%s\n"
	questionPersonMessage  = "次に%sさん，質問お願いします．\n"
	questionEndMessage     = "回答ありがとうございました．"
	personEndMessage       = "これで%sさんの発表時間を終わります．次の発表者は%sさんです．よろしくお願いします．\n"
	meetingEndMessage      = "これで会議を終了します．お疲れ様でした．"
)

func presenOrQuestionEnd(db *gorm.DB, meetingId int, presenterId string, documentId int, isPresenEnd bool) string {
	var endMessage string
	if isPresenEnd {
		endMessage = presenEndMessage
	} else {
		endMessage = questionEndMessage
	}
	isUserId, questionUserId, questionId := selectQuestion(db, meetingId, documentId, presenterId)

	if isUserId { // 質問者を当てる
		questionUserName := getUserName(db, questionUserId)

		return fmt.Sprintf(endMessage+questionPersonMessage, questionUserName)
	} else { // 来ている質問を使う
		questionBody := getQuestionBody(db, questionId)

		return fmt.Sprintf(endMessage+questionBodyAskMessage, questionBody)
	}
}

func personEnd(presenUserId string, nextUserId string) string {
	presenUserName := getUserName(db, presenUserId)
	nextUserName := getUserName(db, nextUserId)

	return fmt.Sprintf(personEndMessage, presenUserName, nextUserName)
}

func meetingEnd() string {
	return meetingEndMessage
}
