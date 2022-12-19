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

package hooks

import "github.com/ainsleyclark/logger/types"

func GetMessage(entry types.Entry, args types.FormatMessageArgs, fmt types.FormatMessageFunc) string {
	// Setup args for formatting the message.
	var message = "" //nolint

	// Use the default format message if none is attached,
	// otherwise call the function that is assigned.
	if fmt == nil {
		message = types.DefaultFormatMessageFn(entry, args)
	} else {
		message = fmt(entry, args)
	}

	return message
}
