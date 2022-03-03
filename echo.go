package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Result struct {
	Result      bool      `json:"result"`
	MeetingName string    `json:"meetingname"`
	StartTime   time.Time `json:"starttime"`
	Presenters  []string  `json:"presenters"`
	DocumentIds []string  `json:"documentids"`
	Scripts     []string  `json:"scripts"`
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
		meetingId, _ := strconv.Atoi(c.FormValue("meetingId"))
		resultJoinMeeting, meetingName, meetingStartTime, presenterNames := joinMeeting(db, c.FormValue("userId"), meetingId)
		test_string := []string{"test"}
		result := &Result{
			Result:      resultJoinMeeting,
			MeetingName: meetingName,
			StartTime:   meetingStartTime,
			Presenters:  presenterNames,
			DocumentIds: test_string,
			Scripts:     test_string,
		}

		return c.JSON(http.StatusOK, result)
	})

	e.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})
}
