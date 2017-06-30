# fasthttplambda

This package allows you to use the `fasthttp` (and `fasthttprouter`) package to create your own service for AWS Lambda.

It also references [eawsy/aws-lambda-go-shim](https://github.com/eawsy/aws-lambda-go-shim) to create the package itself.

## Usage example

```go
package main

import (
	"log"

	"github.com/EtienneBruines/fasthttplambda"
	"github.com/buaazp/fasthttprouter"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/valyala/fasthttp"
)

var router = new(fasthttprouter.Router)

func Handle(event *apigatewayproxyevt.Event, ctx *runtime.Context) (*fasthttplambda.ProxyOutput, error) {
	fasthttplambda.Router = router
	return fasthttplambda.Handle(event, ctx)
}

func main() {
	log.Println("Listening on port 8080")
	fasthttp.ListenAndServe(":8080", router.Handler)
}
```

## Key components

- A `fasthttplambda.Handle` method, which you should call;
- A `fasthttplambda.Router` (of type `*fasthttp.Router`) which you should set beforehand;
	- You can use this Router to define your routes / methods;
- You can use this Router for local development as well, so you don't have to deploy every time.

## How does it work

Currently, it works by catching the request body and headers from the AWS Lambda call, and using them to create an
in-memory `fasthttp.Request` to call to the router we defined.

## Possible optimizations

It's be nice to reference the memory already allocated within the event, instead of copying it into our own `Request`.
However, since the `fasthttp.RequestCtx` is not an `interface`, there's little chance to get it working without forking
the `fasthttp` package.

## Contributing

Any contributions are welcome. This is just a proof-of-concept at this stage. Bug reports / PRs are welcome.
