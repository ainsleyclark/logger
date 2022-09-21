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
	"github.com/ainsleyclark/logger/types"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	report := func(r types.Entry) bool {
		return true
	}
	format := func(entry types.Entry, args types.FormatMessageArgs) string {
		return "test"
	}

	opts := logger.NewOptions().
		Prefix("app").
		Service("api").
		WithWorkplaceNotifier(os.Getenv("WORKPLACE_TOKEN"), os.Getenv("WORKPLACE_THREAD"), report, format)

	err = logger.New(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}

	logger.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error()

	time.Sleep(time.Second * 1)
}
