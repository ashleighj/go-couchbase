package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gocouchbase/pkg/config"
)

const (
	prefixInfo  = " type=\"INFO\"  "
	prefixWarn  = " type=\"WARN\"  "
	prefixError = " type=\"ERROR\" "
)

/**
* For info-level logging.
*
* List of interfaces converted to string to avoid unnecassary [] added by Go
* while being able to add custom prefix
**/
func Info(ctx context.Context, messages ...interface{}) {
	log.Println(getLogString(ctx, prefixInfo, messages))
}

func Infof(ctx context.Context, str string, vars ...interface{}) {
	str = fmt.Sprintf(str, vars...)
	log.Println(getLogString(ctx, prefixInfo, []interface{}{str}))
}

/**
* For logging errors and optional additional information/context messages.
**/
func Error(ctx context.Context, messages ...interface{}) {
	log.Println(getLogString(ctx, prefixError, messages))
}

/*
* For logging non-critical errors as warnings.
 */
func Warn(ctx context.Context, messages ...interface{}) {
	log.Println(getLogString(ctx, prefixWarn, messages))
}

/*
* For logging non-critical errors as warnings.
 */
func Fatal(ctx context.Context, messages ...interface{}) {
	log.Fatal(getLogString(ctx, prefixWarn, messages))
}

// TimeExecution times an executed function
func TimeExecution(ctx context.Context, caller string, t time.Time) {
	elapsed := time.Since(t)
	msg := fmt.Sprintf("function=\"%s\" executedInMs=\"%f\"", caller, float64(elapsed.Nanoseconds())/float64(time.Millisecond))

	Info(ctx, msg)
}

func getLogString(ctx context.Context, prefix string, messages []interface{}) string {
	strMsgs := []string{}
	for _, m := range messages {
		switch m.(type) {
		case error:
			strMsgs = append(strMsgs, m.(error).Error())
		default:
			strMsgs = append(strMsgs, fmt.Sprint(m))
		}
	}
	msgString := strings.Join(strMsgs, ", ")

	logString := prefix
	if prefix == prefixError {
		workDir, _ := os.Getwd()
		_, file, line, _ := runtime.Caller(2)
		file = strings.TrimPrefix(file, workDir+"/")

		errLocation := fmt.Sprintf("\033[31m%s \033[39m", file+":"+strconv.Itoa(line))

		logString += errLocation
	}

	if ctx != nil {
		// get variables to track request
		val := ctx.Value(config.RequestContext)
		//assert data is a map of strings
		data, ok := val.(map[string]string)
		if ok {
			fmStr := "context=(%s=\"%s\", %s=\"%s\", %s=\"%s\", %s=\"%s\", %s=\"%s\", %s=\"%s\")"
			ctxStr := fmt.Sprintf(fmStr,
				config.RequestIDKey, data[config.RequestIDKey],
				config.ClientIDKey, data[config.ClientIDKey],
				config.RequestTimeKey, data[config.RequestTimeKey],
				config.RemoteAddrKey, data[config.RemoteAddrKey],
				config.RequestURIKey, data[config.RequestURIKey],
				config.MethodKey, data[config.MethodKey])

			if msgString != "" {
				pref := fmt.Sprintf("\"%s\"\n\t", msgString)
				ctxStr = pref + ctxStr
			}
			logString += ctxStr
		} else {
			logString += msgString
		}
	} else {
		logString += msgString
	}
	return logString
}
