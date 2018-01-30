//
// +build small

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
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

const createTaskResult = "ID 1000"

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, createTaskResult)
	os.Exit(0)
}

func TestRunCmdRunner(t *testing.T) {
	ExecCommand = fakeExecCommand
	defer func() { ExecCommand = exec.Command }()
	cmdRunner := CmdRunner{}
	_, err := cmdRunner.Run("snaptel", []string{"-u", "localhost", "task", "create", "-t"})

	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}

}
