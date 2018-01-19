package main

import (
	"github.com/EtienneBruines/fasthttplambda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var router fasthttprouter.Router

func HelloWorld(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Hello World!")
}

func HelloUniverse(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Hello Universe!")
}

func main() {
	router.Handle("GET", "/", HelloWorld)
	router.Handle("GET", "/hello-universe", HelloUniverse)

	// Usually, we'd listen on a port and handle it as usual
	//fasthttp.ListenAndServe(":8080", router.Handler)

	// But when compiling it for lambda, we use this instead:
	fasthttplambda.Router = &router
	lambda.Start(fasthttplambda.Handle)
}
