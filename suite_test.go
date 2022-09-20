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
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

// LoggerTestSuite defines the helper used for L
// testing.
type LoggerTestSuite struct {
	suite.Suite
}

// TestLogger asserts testing has begun.
func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

// TearDownTestSuite - Teardown logging after testing.
func (t *LoggerTestSuite) TearDownTestSuite() {
	err := New(context.Background())
	t.NoError(err)
}

// Setup is a helper function for setting up the L
// suite.
func (t *LoggerTestSuite) Setup() *bytes.Buffer {
	buf := &bytes.Buffer{}
	L.SetLevel(logrus.TraceLevel)
	L.SetOutput(buf)
	c := Config{}
	L.SetFormatter(&formatter{
		Config:  c.assignDefaults(),
		Colours: false,
	})
	return buf
}
