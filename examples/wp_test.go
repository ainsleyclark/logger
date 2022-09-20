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

package examples

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/logger"
	"testing"
)

func Test(t *testing.T) {
	opts := logger.NewOptions().
		Service("api").
		WithWorkplaceNotifier("token", "thread")

	err := logger.New(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}

	logger.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error()
}
