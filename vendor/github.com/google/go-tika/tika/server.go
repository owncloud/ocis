/*
Copyright 2017 Google Inc.

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

package tika

import (
	"context"
	"crypto/sha512"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

// Server represents a Tika server. Create a new Server with NewServer,
// start it with Start, and shut it down with the close function returned
// from Start.
// There is no need to create a Server for an already running Tika Server
// since you can pass its URL directly to a Client.
// Additional Java system properties can be added to a Taka Server before
// startup by adding to the JavaProps map
type Server struct {
	jar       string
	url       string // url is derived from port.
	port      string
	cmd       *exec.Cmd
	child     *ChildOptions
	JavaProps map[string]string
}

// ChildOptions represent command line parameters that can be used when Tika is run with the -spawnChild option.
// If a field is less than or equal to 0, the associated flag is not included.
type ChildOptions struct {
	MaxFiles          int
	TaskPulseMillis   int
	TaskTimeoutMillis int
	PingPulseMillis   int
	PingTimeoutMillis int
}

func (co *ChildOptions) args() []string {
	if co == nil {
		return nil
	}
	args := []string{}
	args = append(args, "-spawnChild")
	if co.MaxFiles == -1 || co.MaxFiles > 0 {
		args = append(args, "-maxFiles", strconv.Itoa(co.MaxFiles))
	}
	if co.TaskPulseMillis > 0 {
		args = append(args, "-taskPulseMillis", strconv.Itoa(co.TaskPulseMillis))
	}
	if co.TaskTimeoutMillis > 0 {
		args = append(args, "-taskTimeoutMillis", strconv.Itoa(co.TaskTimeoutMillis))
	}
	if co.PingPulseMillis > 0 {
		args = append(args, "-pingPulseMillis", strconv.Itoa(co.PingPulseMillis))
	}
	if co.PingTimeoutMillis > 0 {
		args = append(args, "-pingTimeoutMillis", strconv.Itoa(co.PingTimeoutMillis))
	}
	return args
}

// URL returns the URL of this Server.
func (s *Server) URL() string {
	return s.url
}

// NewServer creates a new Server. The default port is 9998.
func NewServer(jar, port string) (*Server, error) {
	if jar == "" {
		return nil, fmt.Errorf("no jar file specified")
	}

	// Check if the jar file exists.
	if _, err := os.Stat(jar); os.IsNotExist(err) {
		return nil, fmt.Errorf("jar file %q does not exist", jar)
	}

	if port == "" {
		port = "9998"
	}

	urlString := "http://localhost:" + port
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("invalid port %q: %v", port, err)
	}

	s := &Server{
		jar:       jar,
		port:      port,
		url:       u.String(),
		JavaProps: map[string]string{},
	}

	return s, nil
}

// ChildMode sets up the server to use the -spawnChild option.
// If used, ChildMode must be called before starting the server.
// If you want to turn off the -spawnChild option, call Server.ChildMode(nil).
func (s *Server) ChildMode(ops *ChildOptions) error {
	if s.cmd != nil {
		return fmt.Errorf("server process already started, cannot switch to spawn child mode")
	}
	s.child = ops
	return nil
}

var command = exec.Command

// Start starts the given server. Start will start a new Java process. The
// caller must call Stop() to shut down the process when finished with the
// Server. Start will wait for the server to be available or until ctx is
// cancelled.
func (s *Server) Start(ctx context.Context) error {
	if _, err := os.Stat(s.jar); os.IsNotExist(err) {
		return err
	}

	// Create a slice of Java system properties to be passed to the JVM.
	props := []string{}
	for k, v := range s.JavaProps {
		props = append(props, fmt.Sprintf("-D%s=%q", k, v))
	}

	args := append(append(props, "-jar", s.jar, "-p", s.port), s.child.args()...)
	cmd := command("java", args...)

	if err := cmd.Start(); err != nil {
		return err
	}
	s.cmd = cmd

	if err := s.waitForStart(ctx); err != nil {
		out, readErr := cmd.CombinedOutput()
		if readErr != nil {
			return fmt.Errorf("error reading output: %v", readErr)
		}
		// Report stderr since sometimes the server says why it failed to start.
		return fmt.Errorf("error starting server: %v\nserver stderr:\n\n%s", err, out)
	}
	return nil
}

// waitForServer waits until the given Server is responding to requests or
// ctx is Done().
func (s Server) waitForStart(ctx context.Context) error {
	c := NewClient(nil, s.url)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if _, err := c.Version(ctx); err == nil {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Stop shuts the server down, killing the underlying Java process. Stop
// must be called when finished with the server to avoid leaking the
// Java process. If s has not been started, Stop will panic.
// If not running in a Windows environment, it is recommended to use Shutdown
// for a more graceful shutdown of the Java process.
func (s *Server) Stop() error {
	if err := s.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("could not kill server: %v", err)
	}
	if err := s.cmd.Wait(); err != nil {
		return fmt.Errorf("could not wait for server to finish: %v", err)
	}
	return nil
}

// Shutdown attempts to close the server gracefully before using SIGKILL,
// Stop() uses SIGKILL right away, which causes the kernal to stop the java process instantly.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("could not interrupt server: %v", err)
	}
	errChannel := make(chan error)
	go func() {
		select {
		case errChannel <- s.cmd.Wait():
		case <-ctx.Done():
		}
	}()
	select {
	case err := <-errChannel:
		if err != nil {
			return fmt.Errorf("could not wait for server to finish: %v", err)
		}
	case <-ctx.Done():
		if err := s.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("could not kill server: %v", err)
		}
	}
	return nil
}

func sha512Hash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// A Version represents a Tika Server version.
type Version string

// Supported versions of Tika Server.
const (
	Version119 Version = "1.19"
	Version120 Version = "1.20"
	Version121 Version = "1.21"
)

// Versions is a list of supported versions of Apache Tika.
var Versions = []Version{Version119, Version120, Version121}

var sha512s = map[Version]string{
	Version119: "a9e2b6186cdb9872466d3eda791d0e1cd059da923035940d4b51bb1adc4a356670fde46995725844a2dd500a09f3a5631d0ca5fbc2d61a59e8e0bd95c9dfa6c2",
	Version120: "a7ef35317aba76be8606f9250893efece8b93384e835a18399da18a095b19a15af591e3997828d4ebd3023f21d5efad62a91918610c44e692cfd9bed01d68382",
	Version121: "e705c836b2110530c8d363d05da27f65c4f6c9051b660cefdae0e5113c365dbabed2aa1e4171c8e52dbe4cbaa085e3d8a01a5a731e344942c519b85836da646c",
}

// DownloadServer downloads and validates the given server version,
// saving it at path. DownloadServer returns an error if it could
// not be downloaded/validated.
// It is the caller's responsibility to remove the file when no longer needed.
// If the file already exists and has the correct sha512, DownloadServer will
// do nothing.
func DownloadServer(ctx context.Context, v Version, path string) error {
	hash := sha512s[v]
	if hash == "" {
		return fmt.Errorf("unsupported Tika version: %s", v)
	}
	if got, err := sha512Hash(path); err == nil {
		if got == hash {
			return nil
		}
	}
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	url := fmt.Sprintf("http://search.maven.org/remotecontent?filepath=org/apache/tika/tika-server/%s/tika-server-%s.jar", v, v)
	resp, err := ctxhttp.Get(ctx, nil, url)
	if err != nil {
		return fmt.Errorf("unable to download %q: %v", url, err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("error saving download: %v", err)
	}

	h, err := sha512Hash(path)

	if err != nil {
		return err
	}
	if h != hash {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("invalid sha512: %s: error removing %s: %v", h, path, err)
		}
		return fmt.Errorf("invalid sha512: %s", h)
	}
	return nil
}
