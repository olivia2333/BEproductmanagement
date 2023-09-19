package main

import "github.com/kataras/iris/v12"

func main() {
	// 1. create iris instance
	app := iris.New()
	// 2. error, mvc show error
	app.Logger().SetLevel("debug")
	// 3. register template
	template := iris.HTML("./backend/web/views",
		".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. set template target
	app.HandleDir("/assets", iris.Dir("./backend/web/assets"))
	// on error
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "webpage error"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// 5. register controller

	//6. launch service
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations)
}
