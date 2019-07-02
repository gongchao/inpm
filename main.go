package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/spf13/viper"
	"npm/api/controllers"
	"npm/lib/storage"
	"npm/services"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Create storage
	sg := storage.NewLocalStorage(storage.Configuration{BasePath: viper.GetString("storage.path")})

	app := iris.Default()

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.WriteString(fmt.Sprintf("npm set registry %s://%s", viper.GetString("app.protocol"), viper.GetString("app.host")))
	})

	mvc.New(app.Party("/")).Register(services.NewPackageService(), sg).Handle(new(controllers.PackageController))
	mvc.New(app.Party("/")).Handle(new(controllers.Auth))

	app.Get("*", func(ctx iris.Context) {
		ctx.Application().Logger().Errorf("error", ctx.RequestPath(true))
	})

	err := app.Run(iris.Addr(":80"))
	if err != nil {
		panic(err)
	}
}
