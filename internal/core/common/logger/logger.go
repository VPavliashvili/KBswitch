package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kbswitch/internal/app"
	"kbswitch/internal/core/common/middleware/models"
	"log/slog"
	"net/http"
	"os"
)

const LogIDKey = "logID"

const (
	LevelDebug = slog.LevelDebug
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelFatal: "FATAL",
	LevelTrace: "TRACE",
	LevelDebug: "DEBUG",
}

var lgr *slog.Logger
var lvl = &slog.LevelVar{}

// this function provides missing level string version
// such as FATAL or TRACE
// otherwise slog will print ERROR+4 for FATAL
func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelLabel, exists := LevelNames[level]
		if !exists {
			levelLabel = level.String()
		}

		a.Value = slog.StringValue(levelLabel)
	}

	return a
}

func Init(app app.Application) {
	f, err := os.OpenFile(app.Logging.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}

	var wr io.Writer
	if app.Logging.EnableConsole {
		wr = io.MultiWriter(os.Stdout, f)
	} else {
		wr = f
	}

	lgr = slog.New(slog.NewJSONHandler(wr, &slog.HandlerOptions{
		Level:       lvl,
		ReplaceAttr: ReplaceAttr,
	}))
}

func Info(msg string) {
	lgr.Info(msg)
}

func Error(msg string) {
	lgr.Error(msg)
}

func Warn(msg string) {
	lgr.Warn(msg)
}

func Debug(msg string) {
	lvl.Set(LevelDebug)
	lgr.Debug(msg)
}

func Fatal(msg string) {
	lgr.Log(context.Background(), LevelFatal, msg)
}

func Trace(msg string) {
	lvl.Set(LevelTrace)
	lgr.Log(context.Background(), LevelTrace, msg)
}

func getRequestLog(r *http.Request) models.RequestLog {
	var result models.RequestLog

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	bodyStr := string(body[:])

	result.Body = bodyStr
	result.Method = r.Method

	result.Route = r.URL.String()
	result.Params = r.URL.Query()

	return result
}

func getResponseLog(rww models.ResponseWriterLogWrapper) models.ResponseLog {
	var result models.ResponseLog

	var buf bytes.Buffer
	buf.WriteString(rww.Body.String())

	result.Header = (*rww.W).Header()
	result.Body = buf.String()

	return result
}

// returns json
func GetRequestResponseLog(rww models.ResponseWriterLogWrapper, r *http.Request) string {
	rrl := models.RequestResponseLog{
		Req:        getRequestLog(r),
		Resp:       getResponseLog(rww),
		StatusCode: *(rww.StatusCode),
	}

	bytes, err := json.Marshal(rrl)
	if err != nil {
		panic(fmt.Sprintf("can't decode model for logging, panicking, err: %s", err.Error()))
	}

	return string(bytes[:])
}

func NewResponseWriterWrapper(w http.ResponseWriter) models.ResponseWriterLogWrapper {
	var buf bytes.Buffer
	var statusCode = 200
	return models.ResponseWriterLogWrapper{
		W:          &w,
		Body:       &buf,
		StatusCode: &statusCode,
	}
}

