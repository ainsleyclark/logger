// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package notify

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

// NewFireHook is responsible for creating a workplace client and
// returning a FireHook sending log entire to a workplace
// thread if an error occurred within the system.
func NewFireHook(token string) (mogrus.FireHook, error) {
	const op = "Notify.NewFireHook"
	wp, err := workplace.New(workplace.Config{Token: token})
	if err != nil {
		return nil, errors.NewInvalid(err, "Error creating Workplace Client", op)
	}
	return func(entry mogrus.Entry) {
		go fire(wp, entry)
	}, nil
}

const (
	// Thread is the thread id to send logs to on Workplace.
	Thread = "t_3950977074953346"
)

// fire is a helper that sends messages off to Workplace which
// is called concurrently.
func fire(wp workplace.Notifier, entry mogrus.Entry) {
	// Bail if the error is nil.
	if entry.Error == nil {
		return
	}

	// Bail if the error code is not anything but INTERNAL,
	// we don't want to notify users of invalid or pesky
	// log entries.
	if entry.Error.Code != errors.INTERNAL {
		return
	}

	// Use the Workplace client to send a message via the bot.
	err := wp.Notify(workplace.Transmission{
		Thread:  Thread,
		Message: formatMessage(entry),
	})
	if err != nil {
		log.Println(err.Error()) // We can't use the standard logger as it may cause a loop.
	}
}

// formatMessage prints a formatted message from the log entry to
// a user friendly message.
func formatMessage(entry mogrus.Entry) string {
	buf := bytes.Buffer{}

	// Write Krang & version from the latest build.
	buf.WriteString(fmt.Sprintf("%v Krang v%s\n", emoji.ChartIncreasing, ""))

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
