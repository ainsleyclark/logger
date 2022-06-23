// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

func (t *LoggerTestSuite) TestOptions_AssignDefaults() {
	o := Options{}
	got := o.assignDefaults()
	want := Options{
		Prefix:        DefaultPrefix,
		DefaultStatus: DefaultStatus,
	}
	t.Equal(want, got)
}
