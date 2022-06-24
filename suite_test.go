// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

// LoggerTestSuite defines the helper used for logger
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

// Setup is a helper function for setting up the logger
// suite.
func (t *LoggerTestSuite) Setup() *bytes.Buffer {
	buf := &bytes.Buffer{}
	logger.SetLevel(logrus.TraceLevel)
	logger.SetOutput(buf)
	c := Config{}
	logger.SetFormatter(&formatter{
		Config:  c.assignDefaults(),
		Colours: false,
	})
	return buf
}
