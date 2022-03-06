package main

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Result struct {
	Result bool `json:"result"`
}

type UserSignupRequest struct {
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	UserPassword string `json:"userPassword"`
}

type UserLoginRequest struct {
	UserId       string `json:"userId"`
	UserPassword string `json:"userPassword"`
}

type UserLoginResult struct {
	Result   bool   `json:"result"`
	UserName string `json:"userName"`
}

type CreateMeetingRequest struct {
	MeetingName      string   `json:"meetingName"`
	MeetingStartTime string   `json:"meetingStartTime"`
	Presenters       []string `json:"presenters"`
}

type CreateMeetingResult struct {
	MeetingId        int      `json:"meetingId"`
	MeetingName      string   `json:"meetingName"`
	MeetingStartTime string   `json:"meetingStartTime"`
	Presenters       []string `json:"presenters"`
	DocumentIds      []int    `json:"documentIds"`
}

type JoinMeetingRequest struct {
	UserId    string `json:"userId"`
	MeetingId int    `json:"meetingId"`
}

type JoinMeetingResult struct {
	Result           bool     `json:"result"`
	MeetingName      string   `json:"meetingName"`
	MeetingStartTime string   `json:"meetingStartTime"`
	Presenters       []string `json:"presenters"`
	DocumentIds      []int    `json:"documentIds"`
}

func initRouting(e *echo.Echo, hub *Hub, db *gorm.DB) {

	e.GET("/", func(c echo.Context) error {
		serveHome(c.Response(), c.Request())
		return nil
	})

	e.POST("/user/signup", func(c echo.Context) error {
		request := new(UserSignupRequest)
		err := c.Bind(request)
		if err == nil {
			result := &Result{
				Result: signupUser(db, request.UserId, request.UserName, request.UserPassword),
			}

			return c.JSON(http.StatusOK, result)
		} else {
			return c.JSON(http.StatusBadRequest, &Result{Result: false})
		}
	})

	e.POST("/user/login", func(c echo.Context) error {
		request := new(UserLoginRequest)
		err := c.Bind(request)
		if err == nil {
			resultLogin, userName := loginUser(db, request.UserId, request.UserPassword)
			result := &UserLoginResult{
				Result:   resultLogin,
				UserName: userName,
			}

			return c.JSON(http.StatusOK, result)
		} else {
			return c.JSON(http.StatusBadRequest, &UserLoginResult{Result: false, UserName: ""})
		}
	})

	e.POST("/meeting/join", func(c echo.Context) error {
		request := new(JoinMeetingRequest)
		err := c.Bind(request)
		if err == nil {
			resultJoinMeeting, meetingName, meetingStartTime, presenterNames, documentIds := joinMeeting(db, request.UserId, request.MeetingId)
			layout := "2006/01/02 15:04:05"
			meetingStartTimeString := meetingStartTime.Format(layout)
			result := &JoinMeetingResult{
				Result:           resultJoinMeeting,
				MeetingName:      meetingName,
				MeetingStartTime: meetingStartTimeString,
				Presenters:       presenterNames,
				DocumentIds:      documentIds,
			}
			if result.Result {
				go hub.sendStartMeetingMessage(request.MeetingId, meetingStartTime)
			}
			return c.JSON(http.StatusOK, result)
		} else {
			return c.JSON(http.StatusBadRequest, &Result{Result: false})
		}

	})

	e.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})

	e.POST("/meeting/create", func(c echo.Context) error {
		request := new(CreateMeetingRequest)
		err := c.Bind(request)
		if err == nil {
			meetingId, meetingName, meetingStartTime, presenters, documentIds := createMeeting(db, request.MeetingName, request.MeetingStartTime, request.Presenters)
			result := &CreateMeetingResult{
				MeetingId:        meetingId,
				MeetingName:      meetingName,
				MeetingStartTime: meetingStartTime,
				Presenters:       presenters,
				DocumentIds:      documentIds,
			}

			return c.JSON(http.StatusOK, result)
		} else {
			return c.JSON(http.StatusBadRequest, &Result{Result: false})
		}
	})
}
