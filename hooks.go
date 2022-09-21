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
	"github.com/ainsleyclark/logger/internal/workplace"
	"github.com/ainsleyclark/mogrus"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	// newMogrus is an alias for mogrus.New
	newMogrus = mogrus.New
	// newWP is an alias for notify.NewFireHook
	newWP = workplace.NewHook
)

// addHooks adds all hooks to the logger.
func addHooks(ctx context.Context, cfg *Config) error {
	err := addWorkplaceHook(cfg)
	if err != nil {
		return err
	}

	err = addMogrusHook(ctx, cfg)
	if err != nil {
		return err
	}

	return nil
}

// addWorkplaceHook adds the Workplace hook if
// the thread and token exists.
func addWorkplaceHook(cfg *Config) error {
	if cfg.workplaceThread != "" && cfg.workplaceToken != "" {
		wpHook, err := newWP(workplace.Options{
			Token:        cfg.workplaceToken,
			Thread:       cfg.workplaceThread,
			Service:      cfg.service,
			Version:      cfg.version,
			ShouldReport: cfg.report,
		})
		if err != nil {
			return err
		}
		L.AddHook(wpHook)
	}
	return nil
}

// addMogrusHook adds the Mogrus hook if
// the client exists.
func addMogrusHook(ctx context.Context, cfg *Config) error {
	if cfg.mongoCollection != nil {
		mogrusHook, err := newMogrus(ctx, mogrus.Options{
			Collection: cfg.mongoCollection,
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
		L.AddHook(mogrusHook)
	}
	return nil
}
