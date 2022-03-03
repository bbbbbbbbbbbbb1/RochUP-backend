package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Result struct {
	Result bool `json:"result"`
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

func initRouting(e *echo.Echo, hub *Hub, db *gorm.DB) {

	e.GET("/", func(c echo.Context) error {
		// return c.String(http.StatusOK, "Hello, World!")
		serveHome(c.Response(), c.Request())
		// return c.JSON(http.StatusOK, {"ok": true})
		return nil
	})

	e.GET("/ip", func(c echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(("<h3>あなたのIPアドレスは %s</h3>"), c.RealIP()))
	})

	e.GET("/users/:id", func(c echo.Context) error {
		jsonMap := map[string]string{
			"name": "okutani",
			"hoge": "piyo",
		}
		return c.JSON(http.StatusOK, jsonMap)
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
