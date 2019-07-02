package middlewares

import "github.com/kataras/iris"

func Auth(ctx iris.Context) {
	ctx.Next()
}