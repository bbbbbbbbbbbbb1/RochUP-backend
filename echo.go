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
	UserId       string
	UserName     string
	UserPassword string
}

type UserLoginRequest struct {
	UserId       string
	UserPassword string
}

type CreateMeetingRequest struct {
	MeetingName string
	StartTime   string
	Presenters  []string
}

type CreateMeetingResult struct {
	MeetingId        int
	MeetingName      string
	MeetingStartTime string
	Presenters       []string
	DoncumentIds     []string
	Scripts          []string
}

type JoinMeetingRequest struct {
	UserId    string `json:"userId"`
	MeetingId int    `json:"meetingId"`
}

type JoinMeetingResult struct {
	Result      bool     `json:"result"`
	MeetingName string   `json:"meetingName"`
	StartTime   string   `json:"startTime"`
	Presenters  []string `json:"presenters"`
	DocumentIds []string `json:"documentIds"`
	Scripts     []string `json:"scripts"`
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
		request := new(UserSignupRequest)
		err := c.Bind(request)
		if err == nil {
			result := &Result{
				Result: loginUser(db, request.UserId, request.UserPassword),
			}

			return c.JSON(http.StatusOK, result)
		} else {
			return c.JSON(http.StatusBadRequest, &Result{Result: false})
		}
	})

	e.POST("/meeting/join", func(c echo.Context) error {
		request := new(JoinMeetingRequest)
		err := c.Bind(request)
		if err == nil {
			resultJoinMeeting, meetingName, meetingStartTime, presenterNames := joinMeeting(db, request.UserId, request.MeetingId)
			layout := "2006/01/02 15:04:05"
			meetingStartTimeString := meetingStartTime.Format(layout)
			test_string := []string{"test"}
			result := &JoinMeetingResult{
				Result:      resultJoinMeeting,
				MeetingName: meetingName,
				StartTime:   meetingStartTimeString,
				Presenters:  presenterNames,
				DocumentIds: test_string,
				Scripts:     test_string,
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
			meetingId, meetingName, meetingStartTime, presenters := createMeeting(db, request.MeetingName, request.StartTime, request.Presenters)
			createMeetingResult := &CreateMeetingResult{
				MeetingId:        meetingId,
				MeetingName:      meetingName,
				MeetingStartTime: meetingStartTime,
				Presenters:       presenters,
				DoncumentIds:     []string{},
				Scripts:          []string{},
			}

			return c.JSON(http.StatusOK, createMeetingResult)
		} else {
			return c.JSON(http.StatusBadRequest, &Result{Result: false})
		}
	})
}
