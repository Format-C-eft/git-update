package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/Format-C-eft/git-update/internal/config"
)

const (
	layoutDateFormat = "15:04:05.999"
	gitStatusOk      = "nothing to commit, working tree clean"
)

type (
	resultLog struct {
		dir  string
		logs []logs
	}

	logs struct {
		Time     time.Time
		messages string
	}
)

func (r *resultLog) AddLog(message string, verboseMessage string) {
	r.logs = append(r.logs, logs{
		Time:     time.Now(),
		messages: message,
	})

	if verboseMessage != "" && config.FlagVerbose {
		r.logs = append(r.logs, logs{
			Time:     time.Now(),
			messages: "verbose: " + verboseMessage,
		})
	}
}

func (r *resultLog) String() string {
	result := make([]string, 0, len(r.logs))
	for _, l := range r.logs {
		result = append(result, fmt.Sprintf(
			"%s: %s: %s",
			l.Time.Format(layoutDateFormat),
			r.dir,
			l.messages,
		),
		)
	}

	return strings.Join(result, "\n")
}
