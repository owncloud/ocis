/*
 * Copyright 2017-2019 Kopano and its licensors
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

package kcc

import (
	"net"
	"time"
)

// DefaultUnixDialer is the default Dialer as used by KSS for Unix socket SOAP
// request.
var DefaultUnixDialer = &net.Dialer{
	Timeout: 10 * time.Second,
}

// DefaultUnixMaxConnections is the default maximum number of connections which
// will be created to handle parallel SOAP requests to Unix sockets.
var DefaultUnixMaxConnections = 20
