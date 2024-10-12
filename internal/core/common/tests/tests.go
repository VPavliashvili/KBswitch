package tests

import (
	"encoding/json"
	"fmt"
	"kbswitch/internal/core/common"
	"reflect"
	"testing"
)

var ErrTest = fmt.Errorf("test")

const (
	LogLvlInfo  = "I"
	LogLvlError = "E"
	LogLvlTrace = "T"
)

type FakeLogger struct {
	Logs []string
}

// LogError implements logging.Logger.
func (f *FakeLogger) LogError(msg string) {
	f.Logs = append(f.Logs, LogLvlError)
}

// LogInfo implements logging.Logger.
func (f *FakeLogger) LogInfo(msg string) {
	f.Logs = append(f.Logs, LogLvlInfo)
}

// LogTrace implements logging.Logger.
func (f *FakeLogger) LogTrace(msg string) {
	f.Logs = append(f.Logs, LogLvlTrace)
}

func AssertLogsEqual(method string, t *testing.T, want []string, got []string) {
	if (len(want) != 0 && len(got) != 0) && !reflect.DeepEqual(want, got) || (len(want) != len(got)) {
		t.Errorf("in method %s: log check failed\nexpected %+v\ngot %v", method, want, got)
	}
}

func AssertResultsEqual(method string, t *testing.T, want any, got any) {
	if !reflect.DeepEqual(want, got) {
		w, _ := json.Marshal(want)
		g, _ := json.Marshal(got)
		t.Errorf("in method %s: result check failed\nexpected %s\ngot %s", method, w, g)
	}
}

func AssertHasError(method string, t *testing.T, want error, got error) {
	if want != nil && got == nil {
		t.Errorf("in method %s: expected error is not nil %v, when result returned nil: %v", method, want, got)
	}
}

func AssertErrorsEqual(method string, t *testing.T, want *common.AppError, got *common.AppError) {
	if want == nil && got != nil {
		t.Errorf("in method %s: expected error equals to nil, when error returned: %v", method, got)
	} else if want != nil {
		et := want.Errtype.Error() == got.Errtype.Error()
		er := want.Reason.Error() == got.Reason.Error()
		equals := et && er

		if !equals {
			t.Errorf("in method %s: error check failed\nexpected type: %v\ngot type: %v\nexpected reason: %v\ngot reason: %v",
				method, want.Errtype, got.Errtype, want.Reason, got.Reason)
		}
	}
}
