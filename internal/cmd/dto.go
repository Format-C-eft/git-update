package cmd

import (
	"fmt"
	"strings"
	"time"
)

const (
	LayoutDateFormat = "15:04:05.999"
	gitStatusOk      = "nothing to commit, working tree clean"
)

var (
	flagDir           string
	flagAll           bool
	flagCheckout      bool
	flagFetch         bool
	flagPull          bool
	flagDefaultBranch string

	flagExecuteTimeout *time.Duration
)

type (
	resultLogChan struct {
		dir  string
		logs []logs
	}

	logs struct {
		Time     time.Time
		messages string
	}
)

func (r *resultLogChan) AddLog(message string) {
	r.logs = append(r.logs, logs{
		Time:     time.Now(),
		messages: message,
	})
}

func (r *resultLogChan) String() string {
	result := make([]string, 0, len(r.logs))
	for _, l := range r.logs {
		result = append(result, fmt.Sprintf(
			"%s: %s: %s",
			l.Time.Format(LayoutDateFormat),
			r.dir,
			l.messages,
		),
		)
	}
	return strings.Join(result, "\n")
}
