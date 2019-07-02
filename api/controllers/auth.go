package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
	"github.com/kataras/iris/mvc"
	"time"
)

type Auth struct {
}

func (*Auth) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("PUT", "/-/user/{username}", "Login")
}

type User struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Date     time.Time `json:"date"`
}

func (*Auth) Login(ctx iris.Context) {
	var user User

	if ctx.ReadJSON(&user) != nil {
		ctx.StatusCode(httptest.StatusUnauthorized)
		return
	}

	_, _ = ctx.JSON(iris.Map{"token": "test"})
}
