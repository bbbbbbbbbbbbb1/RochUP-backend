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

type UserSignupRequest struct {
	UserId       string
	UserName     string
	UserPassword string
}

type UserLoginRequest struct {
	UserId       string
	UserPassword string
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

	e.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})
}
