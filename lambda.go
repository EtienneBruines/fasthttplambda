package fasthttplambda

import (
	"encoding/base64"
	"net"
	"net/url"

	"github.com/buaazp/fasthttprouter"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

var Router *fasthttprouter.Router

type ProxyOutput struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Body            string            `json:"body"`
	Headers         map[string]string `json:"headers"`
}

func server(l net.Listener) {
	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}
	err = fasthttp.ServeConn(conn, Router.Handler)
	if err != nil {
		panic(err)
	}
	err = l.Close()
	if err != nil {
		panic(err)
	}
}

func Handle(event *apigatewayproxyevt.Event, ctx *runtime.Context) (*ProxyOutput, error) {
	l := fasthttputil.NewInmemoryListener()
	go server(l)

	var (
		req  = fasthttp.AcquireRequest()
		resp = fasthttp.AcquireResponse()
	)

	uri := url.URL{}
	uri.Path = event.Path
	uri.Host = "localhost"

	vals := url.Values{}
	for k, v := range event.QueryStringParameters {
		vals.Set(k, v)
	}
	uri.RawQuery = vals.Encode()

	req.SetRequestURI(uri.RequestURI())
	req.SetHost("localhost")
	if event.IsBase64Encoded {
		body, err := base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return nil, err
		}
		req.SetBody(body)
	} else {
		req.SetBody([]byte(event.Body))
	}
	for k, v := range event.Headers {
		req.Header.Add(k, v)
	}
	req.Header.SetMethod(event.HTTPMethod)

	client := fasthttp.Client{
		Dial: func(string) (net.Conn, error) { return l.Dial() },
	}

	err := client.Do(req, resp)
	if err != nil {
		panic(err)
	}

	var header = map[string]string{}
	resp.Header.VisitAll(func(k, v []byte) {
		header[string(k)] = string(v)
	})

	var output = &ProxyOutput{
		IsBase64Encoded: false,
		StatusCode:      resp.StatusCode(),
		Body:            string(resp.Body()),
		Headers:         header,
	}

	return output, nil
}
