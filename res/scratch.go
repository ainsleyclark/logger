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

package res

// SetService replaces the service in the configuration and
// creates a new L.
//func SetService(service string) {
//	config.service = service
//	L = logrus.New()
//	color.Greenln(config)
//	err := initialise(context.Background(), config)
//	if err != nil {
//		L.Error(err)
//	}
//}

//func (t *LoggerTestSuite) TestSetService() {
//
//	t.Run("Success", func() {
//		orig := config
//		defer func() {
//			config = orig
//		}()
//		SetService("service")
//		t.Equal("service", config.service)
//	})
//
//	//t.Run("Error", func() {
//	//	buf := t.Setup()
//	//	SetService("")
//	//	color.Greenln(b)
//	//})
//}

//// Bail if the error is nil.
//if entry.Error == nil {
//	return
//}
//
//// Bail if the error code is not anything but INTERNAL,
//// we don't want to notify users of invalid or pesky
//// log entries.
//if entry.Error.Code != errors.INTERNAL {
//	return
//}
