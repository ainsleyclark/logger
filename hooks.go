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

package logger

import (
	"context"
	"github.com/ainsleyclark/logger/internal/hooks/slack"
	"github.com/ainsleyclark/logger/internal/hooks/workplace"
	"github.com/ainsleyclark/logger/types"
	"github.com/ainsleyclark/mogrus"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	// newMogrus is an alias for mogrus.New
	newMogrus = mogrus.New
	// newWP is an alias for notify.NewFireHook
	newWP = workplace.NewHook
	mtx   = sync.Mutex{}
)

type (
	// fireFunc is an alias for firing a logger entry.
	fireFunc func(*logrus.Entry) error
)

// addHooks adds all hooks to the logger.
func addHooks(ctx context.Context, cfg *Config) error {
	d := &defaultHook{
		config: cfg,
	}

	err := d.addWorkplaceHook()
	if err != nil {
		return err
	}

	err = d.addMogrusHook(ctx)
	if err != nil {
		return err
	}

	d.addSlackHook()

	L.AddHook(d)

	return nil
}

// defaultHook is the default hook for processing logger entries.
type defaultHook struct {
	config *Config
	wp     fireFunc
	mogrus fireFunc
	slack  fireFunc
}

// Fire will be called when some logging function is
// called with current hook. It will format log
// entry to string and write it to
// appropriate writer
func (hook *defaultHook) Fire(entry *logrus.Entry) error {
	if entry == nil {
		return nil
	}
	if hook.wp != nil {
		if hook.config.workplace.Report(types.Entry(*entry)) {
			err := hook.wp(entry)
			if err != nil {
				L.WithError(err).Error() // Don't return, still have processing to do.
			}
		}
	}
	if hook.slack != nil {
		if hook.config.slack.Report(types.Entry(*entry)) {
			err := hook.slack(entry)
			if err != nil {
				L.WithError(err).Error() // Don't return, still have processing to do.
			}
		}
	}
	if hook.mogrus != nil {
		if hook.config.mongo.Report(types.Entry(*entry)) {
			go func(fire fireFunc) {
				mtx.Lock()
				err := fire(entry)
				if err != nil {
					L.WithError(err).Error()
				}
				mtx.Unlock()
			}(hook.mogrus)
		}
	}
	return nil
}

// Levels Define on which log levels this hook would
// trigger.
func (hook *defaultHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// addWorkplaceHook adds the Workplace hook if
// the thread and token exists.
func (hook *defaultHook) addWorkplaceHook() error {
	if hook.config.workplace.Thread != "" && hook.config.workplace.Token != "" {
		wpHook, err := newWP(workplace.Options{
			Token:         hook.config.workplace.Token,
			Thread:        hook.config.workplace.Thread,
			FormatMessage: hook.config.workplace.Formatter,
			Args: types.FormatMessageArgs{
				Service: hook.config.service,
				Version: hook.config.version,
				Prefix:  hook.config.prefix,
			},
		})
		if err != nil {
			return err
		}
		hook.wp = wpHook.Fire
	}
	return nil
}

// addMogrusHook adds the Mogrus hook if
// the client exists.
func (hook *defaultHook) addMogrusHook(ctx context.Context) error {
	if hook.config.mongo.Collection != nil {
		mogrusHook, err := newMogrus(ctx, mogrus.Options{
			Collection: hook.config.mongo.Collection,
			ExpirationLevels: mogrus.ExpirationLevels{
				logrus.TraceLevel: time.Hour * 24,
				logrus.DebugLevel: time.Hour * 24,
				logrus.InfoLevel:  time.Hour * 24 * 7,
				logrus.ErrorLevel: time.Hour * 24 * 7 * 4,
				logrus.WarnLevel:  time.Hour * 24 * 7 * 4,
				logrus.PanicLevel: time.Hour * 24 * 7 * 4 * 6,
				logrus.FatalLevel: time.Hour * 24 * 7 * 4 * 6,
			},
		})
		if err != nil {
			return err
		}
		hook.mogrus = mogrusHook.Fire
	}
	return nil
}

func (hook *defaultHook) addSlackHook() {
	if hook.config.slack.Token != "" && hook.config.slack.Channel != "" {
		hook.slack = slack.NewHook(slack.Options{
			Token:         hook.config.slack.Token,
			Channel:       hook.config.slack.Channel,
			FormatMessage: hook.config.slack.Formatter,
			Args: types.FormatMessageArgs{
				Service: hook.config.service,
				Version: hook.config.version,
				Prefix:  hook.config.prefix,
			},
		}).Fire
	}
}
