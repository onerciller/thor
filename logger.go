package thor

import (
	"log"
	"os"
	"time"
)

type logger struct {
	*log.Logger
}

var (
	colorGreen = "\033[32m"
	green      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow     = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset      = string([]byte{27, 91, 48, 109})
)

//Logger for thor
func Logger() HandlerFunc {
	stdlogger := log.New(os.Stdout, "", 0)

	return func(ctx *Context) error {
		// Start timer
		start := time.Now()

		// Process request
		ctx.Next()
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Response.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		stdlogger.Printf("%s[%s]%s %v |%s %3d %s| %12v | %s |%s %-5s %s %s",
			magenta, ctx.Thor.AppName, reset,
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, method, reset,
			ctx.Request.URL.Path,
		)
		return nil
	}

}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	default:
		return white
	}
}
