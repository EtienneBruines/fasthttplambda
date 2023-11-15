package main

import (
	"github.com/Suremeo/fasthttplambda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func HelloWorld(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("Hello World!")
}

func HelloUniverse(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("Hello Universe!")
}

func main() {
	r := router.New()

	r.GET("/", HelloWorld)
	r.GET("/hello-universe", HelloUniverse)

	lambda.Start(fasthttplambda.Handle(r.Handler))
}
