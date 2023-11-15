# fasthttplambda

This package allows you to use the `fasthttp` package to create your own service for AWS Lambda.

It has been updated to use [aws/aws-lambda-go](github.com/aws/aws-lambda-go) instead of the 
[eawsy/aws-lambda-go-shim](https://github.com/eawsy/aws-lambda-go-shim) it previously used. 

## Usage example

See `example/example.go`

## How does it work

Currently, it works by catching the request body and headers from the AWS Lambda call, and using them to create an
in-memory `fasthttp.Request` to call to the router we defined.

## Possible optimizations

It would be nice to reference the memory already allocated within the event, instead of copying it into our own `Request`.
However, since the `fasthttp.RequestCtx` is not an `interface`, there's little chance to get it working without forking
the `fasthttp` package.

## Contributing

Any contributions are welcome. This is just a proof-of-concept at this stage. Bug reports / PRs are welcome.
