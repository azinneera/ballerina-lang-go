// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// interruptExitCode is the conventional 128+SIGINT status for a signalled
// process.
const interruptExitCode = 130

// onInterrupt runs cleanup on SIGINT/SIGTERM and exits, covering the gap that
// deferred cleanup leaves when the process is signalled rather than returning
// normally. The returned stop releases the handler and is safe to call more
// than once.
//
// cleanup must be safe to call concurrently with the caller's own deferred
// cleanup: a signal arriving as the caller returns can run both. At most one
// registration should be live at a time — two handlers both wake on a signal,
// and whichever exits first truncates the other's cleanup.
func onInterrupt(cleanup func()) (stop func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})
	var once sync.Once
	stop = func() {
		once.Do(func() {
			signal.Stop(ch)
			close(done)
		})
	}
	go func() {
		select {
		case <-ch:
			// Restore the previous disposition before the (potentially slow)
			// cleanup, so a second signal is acted on immediately instead of
			// queueing behind it.
			signal.Reset(os.Interrupt, syscall.SIGTERM)
			cleanup()
			os.Exit(interruptExitCode)
		case <-done:
		}
	}()
	return stop
}
