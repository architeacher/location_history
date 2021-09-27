package infrastructure

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type (
	// RequestLogMessage wrapper for log message data.
	RequestLogMessage struct {
		// Log format CommonLogFormat, or
		Format string
		// Request to be logged.
		Request *http.Request
		// StatusCode of the response.
		StatusCode int
		// BodyLength in bytes.
		BodyLength int
	}

	Logger interface {
		Info(...interface{}) error
		Infof(string, ...interface{}) error
		Error(...interface{}) error
		Errorf(string, ...interface{}) error
	}

	LoggingQueue chan string
)

const (
	// baseFormat have the base for every logging format.
	// Usually used for request that sets only header, and does not return any body.
	// e.g 125.125.125.125 - akamal [10/Oct/1999:21:15:05 +0500] "GET /index.html?q={search} HTTP/1.0"
	baseFormat = "%s %s %s [%s] \"%s %s %s\""
	// CommonLogFormat host rfc931 username date:time "method request protocol" status code
	// e.g. 125.125.125.125 - akamal [10/Oct/1999:21:15:05 +0500] "GET /index.html?q={search} HTTP/1.0" 200 1043
	// For more details: https://en.wikipedia.org/wiki/Common_Log_Format
	CommonLogFormat string = baseFormat + " %d %d"
	// CombinedLogFormat host rfc931 username date:time "method request protocol" status code bytes "referrer" "user_agent" "cookie"
	// e.g. 125.125.125.125 - akamal [10/Oct/1999:21:15:05 +0500] "GET /index.html?q={search} HTTP/1.0" 200 1043
	// "http://www.ibm.com/" "Mozilla/4.05 [en] (WinNT; I)" "USERID=CustomerA;IMPID=01234"
	//
	// For more details: http://publib.boulder.ibm.com/tividd/td/ITWSA/ITWSA_info45/en_US/HTML/guide/c-logs.html#combined
	CombinedLogFormat = CommonLogFormat + " %s %s %s"

	dateFormat string = "01/Jan/2006:15:04:05 -0700"
)

var (
	ReqLogQueue = make(chan *RequestLogMessage, 10)
	LogQueue    = make(LoggingQueue, 10)
)

// LoggingMiddleware returns an HTTP handler.
func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LogRequest(ReqLogQueue, &RequestLogMessage{
			Format:  baseFormat,
			Request: r,
		})
		h.ServeHTTP(w, r)
	})
}

func buildValues(message *RequestLogMessage) []interface{} {
	r := message.Request
	values := []interface{}{
		r.RemoteAddr,
		"-",
		"-",
		time.Now().Format(dateFormat),
		r.Method,
		r.URL.RequestURI(),
		r.Proto,
	}

	switch message.Format {
	case CombinedLogFormat, CommonLogFormat:
		values = append(values, message.StatusCode, message.BodyLength)
	}

	if message.Format != CombinedLogFormat {
		return values
	}

	referer := "-"
	if r.Referer() != "" {
		referer = fmt.Sprintf("\"%s\"", r.Referer())
	}

	userAgent := "-"
	if r.UserAgent() != "" {
		userAgent = fmt.Sprintf("\"%s\"", r.UserAgent())
	}

	cookiesRaw := "-"
	if len(r.Cookies()) > 0 {
		var cookiesValues []string
		for _, cookie := range r.Cookies() {
			cookiesValues = append(cookiesValues, cookie.String())
		}

		cookiesRaw = fmt.Sprintf("\"%s\"", strings.Join(cookiesValues, ";"))
	}

	return append(values, referer, userAgent, cookiesRaw)
}

// LogRequest logs requests to the CLI.
func LogRequest(reqLogQueue chan *RequestLogMessage, message *RequestLogMessage) {
	reqLogQueue <- message
}

func MonitorRequestLogMessages(logQueue chan *RequestLogMessage) {
	for message := range logQueue {
		values := buildValues(message)
		logrus.Infof(message.Format, values...)
	}
}

func MonitorLogMessages(logQueue LoggingQueue) {
	for message := range logQueue {
		logrus.Infof(message)
	}
}

func MonitorErrors(errChan chan error) {
	for err := range errChan {
		if err != nil && err != io.EOF {
			logrus.Error(err)
		}
	}
}
