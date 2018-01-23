/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2018 Intel Corporation

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

package runner

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

//ExecCommand variable
var ExecCommand = exec.Command

//CmdRunner struct
type CmdRunner struct{}

// New instance of CmdRunner
func New() *CmdRunner {
	return &CmdRunner{}
}

//Run runs the command
func (c *CmdRunner) Run(cmd string, args []string) (io.Reader, error) {
	command := ExecCommand(cmd, args...)
	resCh := make(chan []byte)
	errCh := make(chan error)
	go func() {
		out, err := command.CombinedOutput()
		if err != nil {
			errCh <- err
		}
		resCh <- out
	}()
	timer := time.After(2 * time.Second)
	select {
	case err := <-errCh:
		return nil, err
	case res := <-resCh:
		return bytes.NewReader(res), nil
	case <-timer:
		return nil, fmt.Errorf("time out (cmd:%v args:%v)", cmd, args)
	}
}
