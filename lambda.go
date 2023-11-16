package fasthttplambda

import (
	"encoding/base64"
	"net"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type Request struct {
	// Lambda function url methods
	Version        string                                 `json:"version"` // Version is expected to be `"2.0"`
	RawPath        string                                 `json:"rawPath"`
	RawQueryString string                                 `json:"rawQueryString"`
	RequestContext events.LambdaFunctionURLRequestContext `json:"requestContext"`

	// Api Gateway methods
	Path       string `json:"path"` // The url path for the caller
	HTTPMethod string `json:"httpMethod"`

	// Shared methods

	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
	Cookies               []string          `json:"cookies,omitempty"`
	Headers               map[string]string `json:"headers"`
	Body                  string            `json:"body,omitempty"`
	IsBase64Encoded       bool              `json:"isBase64Encoded"`
}

type Response struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	Cookies         []string          `json:"cookies"`
}

func server(l net.Listener, handler fasthttp.RequestHandler) {
	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}
	err = fasthttp.ServeConn(conn, handler)
	if err != nil {
		panic(err)
	}
	err = l.Close()
	if err != nil {
		panic(err)
	}
}

func Handle(handler fasthttp.RequestHandler) func(event Request) (Response, error) {
	return func(event Request) (Response, error) {
		//func Handle(event *apigatewayproxyevt.Event, ctx *runtime.Context) (*ProxyOutput, error) {
		l := fasthttputil.NewInmemoryListener()
		go server(l, handler)

		var (
			req  = fasthttp.AcquireRequest()
			resp = fasthttp.AcquireResponse()
		)

		uri := url.URL{}

		uri.Path = event.Path
		if uri.Path == "" {
			uri.Path = event.RequestContext.HTTP.Path
		}

		uri.Host = "localhost"

		uri.RawQuery = event.RawQueryString
		if uri.RawQuery == "" {
			vals := url.Values{}
			for k, v := range event.QueryStringParameters {
				vals.Set(k, v)
			}
			uri.RawQuery = vals.Encode()
		}

		req.SetRequestURI(uri.RequestURI())
		req.SetHost("localhost")
		if event.IsBase64Encoded {
			body, err := base64.StdEncoding.DecodeString(event.Body)
			if err != nil {
				return Response{}, err
			}
			req.SetBody(body)
		} else {
			req.SetBody([]byte(event.Body))
		}

		for k, v := range event.Headers {
			req.Header.Add(k, v)
		}

		if event.HTTPMethod == "" {
			req.Header.SetMethod(event.RequestContext.HTTP.Method)
		} else {
			req.Header.SetMethod(event.HTTPMethod)
		}

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

		var output = Response{
			IsBase64Encoded: true,
			StatusCode:      resp.StatusCode(),
			Body:            base64.RawStdEncoding.EncodeToString(resp.Body()),
			Headers:         header,
		}

		return output, nil
	}
}
