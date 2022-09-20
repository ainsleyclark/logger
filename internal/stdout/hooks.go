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

package stdout

import (
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"io"
)

// Hook is a hook that writes logs of specified
// LogLevels to specified Writer.
type Hook struct {
	// The io.Writer, this can be stdout or stderr.
	Writer io.Writer
	// The slice of log levels the writer can too.
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is
// called with current hook. It will format log
// entry to string and write it to
// appropriate writer
func (hook *Hook) Fire(entry *logrus.Entry) error {
	const op = "Logger.Hook.Fire"

	line, err := entry.String()
	if err != nil {
		return &errors.Error{Code: errors.INTERNAL, Message: "Error obtaining the entry string", Operation: op, Err: err}
	}

	_, err = hook.Writer.Write([]byte(line))
	if err != nil {
		return &errors.Error{Code: errors.INTERNAL, Message: "Error writing entry to io.Writer", Operation: op, Err: err}
	}

	return nil
}

// Levels Define on which log levels this hook would
// trigger.
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}
