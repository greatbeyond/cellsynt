// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Written by David Högborg <d@greatbeyond.se>, 2016
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cellsynt

import (
	"fmt"
	"regexp"
	"runtime"
)

func ternaryStr(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func caller() string {
	_, file, line, _ := runtime.Caller(1)
	matcher := regexp.MustCompile("^(.*)/(.*?)\\.go$")
	matches := matcher.FindAllStringSubmatch(file, -1)
	msg := fmt.Sprintf(" [cellsynt/%s.go:%d]", matches[0][2], line)

	return msg
}
