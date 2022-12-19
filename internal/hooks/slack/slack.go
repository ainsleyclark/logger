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

package slack

import (
	"github.com/ainsleyclark/logger/internal/hooks"
	"github.com/ainsleyclark/logger/types"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"log"
)

// NewHook creates a new Workplace hook.
// Returns an error if the client could not be created.
func NewHook(opts Options) *Hook {
	return &Hook{
		sendFunc:  slack.New(opts.Token).PostMessage,
		options:   opts,
		LogLevels: logrus.AllLevels,
	}
}

type (
	// Hook represents the Slack hook notifier
	// for log entries.
	Hook struct {
		sendFunc  sendSlackFunc
		options   Options
		LogLevels []logrus.Level
	}
	// Options defines the configuration needed to fire logs
	// via the Workplace API.
	Options struct {
		Token         string
		Channel       string
		Args          types.FormatMessageArgs
		FormatMessage types.FormatMessageFunc
	}
	// sendSlackFunc is the function used for sending to
	// a slack channel.
	sendSlackFunc func(channelID string, options ...slack.MsgOption) (string, string, error)
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
	// Create the Slack attachment that we will send to the channel
	attachment := slack.Attachment{
		Pretext: "Logger",
		Text:    hooks.GetMessage(entry, hook.options.Args, hook.options.FormatMessage),
	}
	// Use the Slack client to send a message via the bot.
	_, _, err := hook.sendFunc(hook.options.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		log.Println(err.Error()) // We can't use the standard logger as it may cause a loop.
	}
}

// Levels Define on which log levels this hook would
// trigger.
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}
