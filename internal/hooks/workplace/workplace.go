// Copyright 2022 Ainsley Clark. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workplace

import (
	"github.com/ainsleyclark/logger/internal/hooks"
	"github.com/ainsleyclark/logger/types"
	"github.com/ainsleyclark/workplace"
	"github.com/sirupsen/logrus"
	"log"
)

// NewHook creates a new Workplace hook.
// Returns an error if the client could not be created.
func NewHook(opts Options) (*Hook, error) {
	wp, err := workplace.New(workplace.Config{Token: opts.Token})
	if err != nil {
		return nil, err
	}
	return &Hook{
		wp:        wp,
		options:   opts,
		LogLevels: logrus.AllLevels,
	}, nil
}

type (
	// Hook represents the workplace hook notifier
	// for log entries.
	Hook struct {
		wp        workplace.Notifier
		options   Options
		LogLevels []logrus.Level
	}
	// Options defines the configuration needed to fire logs
	// via the Workplace API.
	Options struct {
		Token         string
		Thread        string
		Args          types.FormatMessageArgs
		FormatMessage types.FormatMessageFunc
	}
)

// Fire will be called when some logging function is
// called with current hook. It will format log
// entry to string and write it to
// appropriate writer
func (hook *Hook) Fire(entry *logrus.Entry) error {
	go hook.process(types.Entry(*entry)) // This is already check for nil pointer.
	return nil
}

func (hook *Hook) process(entry types.Entry) {
	// Use the Workplace client to send a message via the bot.
	err := hook.wp.Notify(workplace.Transmission{
		Thread:  hook.options.Thread,
		Message: hooks.GetMessage(entry, hook.options.Args, hook.options.FormatMessage),
	})
	if err != nil {
		log.Println(err.Error()) // We can't use the standard logger as it may cause a loop.
	}
}

// Levels Define on which log levels this hook would
// trigger.
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}
