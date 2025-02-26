/*
Copyright 2017 Google Inc. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package tika provides a client and server for downloading, starting, and using Apache Tika's (http://tika.apache.org) Server.

Start with basic imports:

	import "github.com/google/go-tika/tika"

You will need a running Server to make API calls to. So, if you don't
have a server that is already running, and you don't have the Server
JAR already downloaded, you can download one. The caller is responsible
for removing the file when no longer needed.

Version is a custom type, and should be passed as such. There are constants in the code for these.
The following example downloads version 1.21 to the named JAR in the
current working directory.

	err := tika.DownloadServer(context.Background(), tika.Version121, "tika-server-1.21.jar")
	if err != nil {
		log.Fatal(err)
	}

If you don't have a running Tika Server, you can start one.

	// Optionally pass a port as the second argument.
	s, err := tika.NewServer("tika-server-1.21.jar", "")
	if err != nil {
		log.Fatal(err)
	}
	err := s.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer s.Stop()

To parse the contents of a file (or any io.Reader), you will need to open the io.Reader,
create a client, and call client.Parse.

	// import "os"
	f, err := os.Open("path/to/file")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	client := tika.NewClient(nil, s.URL())
	body, err := client.Parse(context.Background(), f)

If you pass an *http.Client to tika.NewClient, it will be used for all requests.

Some functions return a custom type, like Parsers(), Detectors(), and
MIMETypes(). Use these to see what features are supported by the current
Tika server.
*/
package tika
