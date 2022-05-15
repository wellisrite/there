package middlewares

import (
	"fmt"
	. "github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/utils/color"
	"log"
	"net/http"
	"strconv"
	"time"
)

type responseWriterWrapper struct {
	writer        http.ResponseWriter
	writtenHeader int
	writtenBytes  *[]byte
}

func (r *responseWriterWrapper) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriterWrapper) Write(bytes []byte) (int, error) {
	r.writtenBytes = &bytes
	return r.writer.Write(bytes)
}

func (r *responseWriterWrapper) WriteHeader(statusCode int) {
	r.writtenHeader = statusCode
	r.writer.WriteHeader(statusCode)
}

type LoggerConfiguration struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

func Logger(configuration ...LoggerConfiguration) func(request Request, next Response) Response {

	config := &LoggerConfiguration{
		InfoLogger:  log.Default(),
		ErrorLogger: log.Default(),
	}

	if len(configuration) >= 1 {
		config = &configuration[0]
	}

	return func(request Request, next Response) Response {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := &responseWriterWrapper{
				writer:        w,
				writtenHeader: StatusOK,
				writtenBytes:  &[]byte{},
			}
			start := time.Now()
			defer func() {
				code := ww.writtenHeader
				diff := time.Now().Sub(start)
				toLog := color.Blue(r.Method+" "+r.URL.Path) + " resulted in " + statusCodeToColoredString(code) + " (" + StatusText(code) + ") after " + fmt.Sprint(diff)

				if code == StatusInternalServerError {
					config.ErrorLogger.Println(toLog+":", string(*ww.writtenBytes))
				} else {
					config.InfoLogger.Println(toLog)
				}
			}()
			next.ServeHTTP(ww, r)
		}
		return HttpResponseFunc(fn)
	}
}

func statusCodeToColoredString(code int) string {
	c := code - (code % 100)
	cs := strconv.Itoa(code)
	switch c {
	case 100:
		return color.Yellow(cs)
	case 200:
		return color.Green(cs)
	case 300:
		return color.Yellow(cs)
	case 400:
		return color.Magenta(cs)
	case 500:
		return color.Red(cs)
	default:
		return color.Reset(cs)
	}
}
