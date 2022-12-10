package main

import "github.com/kataras/iris/v12"

func main() {

	app := iris.New()

	app.Get("/", func(ctx iris.Context) {

		ctx.WriteString("Hello world with air")

	})

	app.Listen(":8080")

}
