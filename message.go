package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	presenEndMessage       = "発表ありがとうございました．\n"
	questionBodyAskMessage = "匿名質問です．%dページについての質問です．%s\n"
	questionPersonMessage  = "次に%sさん，質問お願いします．\n"
	questionEndMessage     = "回答ありがとうございました．\n"
	personEndMessage       = "これで%sさんの発表時間を終わります．次の発表者は%sさんです．よろしくお願いします．\n"
	meetingStartMessage    = "これから会議を開始します．\n"
	meetingEndMessage      = "これで会議を終了します．お疲れ様でした．\n"
)

func presenOrQuestionEnd(db *gorm.DB, meetingId int, presenterId string, isPresenEnd bool) (msg, qUserId string, qId int) {
	var (
		endMessage     string
		pickQuestioner bool
		dPage          int
	)
	if isPresenEnd {
		endMessage = presenEndMessage
	} else {
		endMessage = questionEndMessage
	}
	pickQuestioner, qUserId, qId = selectQuestion(db, meetingId, getDocumentId(db, presenterId, meetingId), presenterId)

	if pickQuestioner { // 質問者を当てる
		qUserName := getUserName(db, qUserId)
		msg = fmt.Sprintf(endMessage+questionPersonMessage, qUserName)
		return msg, qUserId, qId
	} else { // 来ている質問を使う
		var qBody string
		qBody, dPage = getQuestionBody(db, qId)
		msg = fmt.Sprintf(endMessage+questionBodyAskMessage, dPage, qBody)
		return msg, "", qId
	}
}

func personEnd(presenUserId string, nextUserId string, meetingId int) string {
	presenUserName := getUserName(db, presenUserId)
	nextUserName := getUserName(db, nextUserId)

	return fmt.Sprintf(personEndMessage, presenUserName, nextUserName)
}

func meetingStart() string {
	return meetingStartMessage
}

func meetingEnd() string {
	return meetingEndMessage
}
