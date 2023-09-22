package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
	"seckill-product/backend/web/controllers"
	"seckill-product/common"
	"seckill-product/repositories"
	"seckill-product/services"
)

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
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5. register controller
	productRepository := repositories.NewProductManager(
		"product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	//6. launch service
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations)
}
