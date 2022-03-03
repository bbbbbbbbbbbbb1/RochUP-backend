package main

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Result struct {
	Result bool `json:"result"`
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
		result := &Result{
			Result: signupUser(db, c.FormValue("userId"), c.FormValue("userName"), c.FormValue("userPassword")),
		}

		return c.JSON(http.StatusOK, result)
	})

	e.POST("/user/login", func(c echo.Context) error {
		result := &Result{
			Result: loginUser(db, c.FormValue("userId"), c.FormValue("userPassword")),
		}

		return c.JSON(http.StatusOK, result)
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
}
