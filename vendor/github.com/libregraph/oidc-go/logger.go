/*
 * Copyright 2019 Kopano
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oidc

type logger interface {
	Printf(string, ...interface{})
}

type noopLogger struct {
}

func (log *noopLogger) Printf(string, ...interface{}) {
}

// DefaultLogger is the logger used by this library if no other is explicitly
// specified.
var DefaultLogger logger = &noopLogger{}
