// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workplace

import (
	"bytes"
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/mogrus"
	"github.com/ainsleyclark/workplace"
	"github.com/enescakir/emoji"
	"github.com/sirupsen/logrus"
	"log"
)

// NewWorkplaceHook creates a new Workplace hook.
// Returns an error if the client could not be created.
func NewWorkplaceHook(opts Options) (*WorkplaceHook, error) {
	wp, err := workplace.New(workplace.Config{Token: opts.Token})
	if err != nil {
		return nil, err
	}
	return &WorkplaceHook{
		wp:        wp,
		options:   opts,
		LogLevels: logrus.AllLevels,
	}, nil
}

type (
	// WorkplaceHook represents the workplace hook notifier
	// for log entries.
	WorkplaceHook struct {
		wp        workplace.Notifier
		options   Options
		LogLevels []logrus.Level
	}
	// Options defines the configuration needed to fire logs
	// via the Workplace API.
	Options struct {
		Token   string
		Thread  string
		Service string
		Version string
	}
)

// Fire will be called when some logging function is
// called with current hook. It will format log
// entry to string and write it to
// appropriate writer
func (hook *WorkplaceHook) Fire(entry *logrus.Entry) error {
	go func() {
		formatted := mogrus.ToEntry(entry)

		// Bail if the error is nil.
		if formatted.Error == nil {
			return
		}

		// Bail if the error code is not anything but INTERNAL,
		// we don't want to notify users of invalid or pesky
		// log entries.
		if formatted.Error.Code != errors.INTERNAL {
			return
		}

		// Use the Workplace client to send a message via the bot.
		err := hook.wp.Notify(workplace.Transmission{
			Thread:  hook.options.Thread,
			Message: hook.formatMessage(formatted),
		})
		if err != nil {
			log.Println(err.Error()) // We can't use the standard logger as it may cause a loop.
		}
	}()

	return nil
}

// formatMessage prints a formatted message from the log entry to
// a user friendly message.
func (hook *WorkplaceHook) formatMessage(entry mogrus.Entry) string {
	buf := bytes.Buffer{}

	// Write Krang & version from the latest build.
	buf.WriteString(fmt.Sprintf("%v %s v%s\n", hook.options.Service, emoji.ChartIncreasing, hook.options.Version))

	// Write intro text.
	buf.WriteString("\U0001FAE0 Error detected in Krang, please see the information below for more details.\n\n")

	// Write log
	buf.WriteString(fmt.Sprintf("%v Level: %s\n", emoji.RightArrow, entry.Level))
	buf.WriteString(fmt.Sprintf("%v Time: %s\n", emoji.RightArrow, entry.Time.String()))
	if entry.Message != "" {
		buf.WriteString(fmt.Sprintf("%v Message: %s\n", emoji.RightArrow, entry.Message))
	}

	// Print out the Entries error.
	buf.WriteString(fmt.Sprintf("%v Code: %s\n", emoji.RightArrow, entry.Error.Code))
	buf.WriteString(fmt.Sprintf("%v Message: %s\n", emoji.RightArrow, entry.Error.Message))
	buf.WriteString(fmt.Sprintf("%v Operation: %s\n", emoji.RightArrow, entry.Error.Operation))
	buf.WriteString(fmt.Sprintf("%v Error: %s\n", emoji.RightArrow, entry.Error.Err))
	buf.WriteString(fmt.Sprintf("%v Fileline: %s\n\n", emoji.RightArrow, entry.Error.FileLine))

	// Print out associated data.
	if len(entry.Data) > 0 {
		buf.WriteString("Log entries:\n")
		for k, v := range entry.Data {
			if k == logrus.ErrorKey {
				continue
			}
			buf.WriteString(fmt.Sprintf("%s: %v\n", k, v))
		}
	}

	return buf.String()
}

// Levels Define on which log levels this hook would
// trigger.
func (hook *WorkplaceHook) Levels() []logrus.Level {
	return hook.LogLevels
}
