package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	rwriter "bitbucket.movista.ru/maas/maasapi/internal/response_writer"
	"github.com/labstack/echo/v4"
)

const LogTime = "2006-01-02 15:04:05.000"
const ServiceName = "maasapi"
const DefaultRequestID = "maasapi-internal"

type JsonLogger struct {
}

type LogFields map[string]interface{}

func NewJsonLogger() *JsonLogger {
	return &JsonLogger{}
}

func (jl *JsonLogger) Print(in LogFields) {
	in["service"] = ServiceName
	in["time"] = time.Now().Format(LogTime)
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	// enc.SetIndent("", "    ")
	enc.Encode(in)
}

func (jl *JsonLogger) Info(in LogFields) {
	in["level"] = "INFORMATION"
	jl.Print(in)
}

func (jl *JsonLogger) Error(in LogFields) {
	in["level"] = "ERROR"
	jl.Print(in)
}

func tryUnmarshal(body []byte) interface{} {
	var b interface{}
	err := json.Unmarshal(body, &b)
	if err != nil {
		// not a json body
		b = string(body)
	}
	return b
}

func (jl *JsonLogger) LogSubrequest(req *http.Request, reqID string, spanReqID string,
	token interface{}) error {

	var err error
	body := make([]byte, 0)
	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	if reqID == "" {
		reqID = DefaultRequestID
	}
	jl.Info(LogFields{
		"type":            "SUBREQUEST",
		"request_id":      reqID,
		"subrequest_id":   spanReqID,
		"subrequest_type": req.Method,
		"auth":            token,
		"url":             req.URL.String(),
		"payload":         tryUnmarshal(body),
		"query":           req.URL.Query(),
	})

	return nil
}

func (jl *JsonLogger) LogSubrequestResponse(req *http.Request, resp *http.Response, reqID string, spanReqID string, start time.Time, body []byte) {
	code := resp.StatusCode
	logFunc := jl.Info
	if code >= 400 {
		logFunc = jl.Error
	}
	if reqID == "" {
		reqID = DefaultRequestID
	}
	logFunc(LogFields{
		"type":              "SUBREQUEST_RESPONSE",
		"request_id":        reqID,
		"subrequest_id":     spanReqID,
		"url":               req.URL.String(),
		"status":            code,
		"estimated_time_ms": time.Since(start).Milliseconds(),
		"body":              tryUnmarshal(body),
	})
}

func (jl *JsonLogger) LogSubrequestError(req *http.Request, reqID, spanReqID string, start time.Time, errs []error) {
	if reqID == "" {
		reqID = DefaultRequestID
	}
	errsStr := make([]string, len(errs), len(errs))
	for i, err := range errs {
		errsStr[i] = err.Error()
	}
	jl.Error(LogFields{
		"type":              "SUBREQUEST_RESPONSE",
		"request_id":        reqID,
		"subrequest_id":     spanReqID,
		"url":               req.URL.String(),
		"estimated_time_ms": time.Since(start).Milliseconds(),
		"errors":            errsStr,
	})
}

func getFullURL(req *http.Request) string {
	host := req.Header.Get("Host")
	return host + req.URL.String()
}

func (jl *JsonLogger) LogPanic(c echo.Context, err error, stack []byte) {
	st := string(stack)
	st = strings.ReplaceAll(st, "\t", "")

	reqID := c.Get("requestId").(string)
	jl.Error(LogFields{
		"type":       "EXCEPTION",
		"request_id": reqID,
		"exception":  fmt.Sprintf("%v", err),
		"stacktrace": strings.Split(st, "\n"),
	})
}

func (jl *JsonLogger) LogRequest(c echo.Context, token interface{}) {
	url := getFullURL(c.Request())
	method := path.Base(c.Request().URL.String())
	body := c.Get("reqBody").([]byte)
	reqID := c.Get("requestId").(string)
	jl.Info(LogFields{
		"type":        "REQUEST",
		"url":         url,
		"method":      method,
		"method_type": c.Request().Method,
		"payload":     tryUnmarshal(body),
		"auth":        token,
		"request_id":  reqID,
		"query":       c.Request().URL.Query(),
	})
}

func (jl *JsonLogger) LogResponse(c echo.Context, err error) {
	url := c.Request().URL.String()
	method := path.Base(url)
	body := c.Response().Writer.(*rwriter.ResponseWriter).Body.Bytes()
	start := c.Get("reqStart").(time.Time)
	reqID := c.Get("requestId").(string)
	code := c.Response().Status
	logFunc := jl.Info
	if code >= 400 || err != nil {
		logFunc = jl.Error
	}

	fields := LogFields{
		"type":              "RESPONSE",
		"url":               url,
		"method":            method,
		"method_type":       c.Request().Method,
		"body":              tryUnmarshal(body),
		"request_id":        reqID,
		"response_code":     code,
		"method_elapsed_ms": time.Since(start).Milliseconds(),
	}
	if err != nil {
		fields["error"] = fmt.Sprintf("%v", err)
	}

	logFunc(fields)
}
