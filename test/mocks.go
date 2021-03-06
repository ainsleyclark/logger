// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"github.com/ainsleyclark/workplace"
	"github.com/sirupsen/logrus"
)

type Notifier interface {
	workplace.Notifier
}

type Hook interface {
	logrus.Hook
}
