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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/intelsdi-x/snap-deploy/runner"

	. "github.com/smartystreets/goconvey/convey"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func fakeRemoveCommand(file string) error {
	return nil
}

const createTaskResult = "ID 1000"

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, createTaskResult)
	os.Exit(0)
}

func TestCli(t *testing.T) {
	// Using package functions
	Convey("Should retrieve snap url", t, func() {
		config := ConfigAPI{}
		config.SnapPort = "8181"
		out := getSnapURL(config)
		So(out, ShouldResemble, "http://localhost:8181")
	})
	Convey("Should retrieve tags", t, func() {
		config := ConfigAPI{}
		config.Tags = "one:two,three:four"
		out := unpackTags(config)
		tagMap := map[string]string{"one": "two", "three": "four"}
		So(out, ShouldResemble, tagMap)
	})
	Convey("Should generate task", t, func() {
		config := ConfigAPI{}
		config.Tags = "one:two,three:four"
		config.SnapPort = "8181"
		config.DbHost = "localhost"
		config.Metrics = "/intel"
		config.DbPassword = "snap"
		config.DbUser = "snap"
		config.DbDatabase = "snap"
		out, err := generateTask(config)
		expectedOutput := "{\"version\":1,\"schedule\":{\"type\":\"simple\",\"interval\":\"\"},\"workflow\":{\"collect\":{\"metrics\":{\"/intel\":{}},\"tags\":{\"/intel\":{\"one\":\"two\",\"three\":\"four\"}},\"publish\":[{\"plugin_name\":\"influxdb\",\"config\":{\"host\":\"localhost\",\"port\":8086,\"database\":\"snap\",\"user\":\"snap\",\"password\":\"snap\"}}],\"config\":{\"/intel/libvirt\":{\"nova\":true}}}}}"
		So(string(out), ShouldResemble, expectedOutput)
		So(err, ShouldBeNil)
	})

	Convey("Should create task", t, func() {
		runner.ExecCommand = fakeExecCommand
		config := ConfigAPI{}
		config.Tags = "one:two,three:four"
		config.SnapPort = "8181"
		config.DbHost = "localhost"
		config.Metrics = "/intel"
		config.DbPassword = "snap"
		config.DbUser = "snap"
		config.DbDatabase = "snap"
		out, err := createTaskCli(config)
		expectedOutput := "1000"
		fmt.Println(out)
		So(string(out), ShouldResemble, expectedOutput)
		So(err, ShouldBeNil)
	})

	Convey("Should download plugin", t, func() {

		config := ConfigAPI{}
		config.plugins = "collector-psutil"
		config.snapLocation = "/tmp"
		removeCommand = fakeRemoveCommand
		download(config)
		fileInfo, err := os.Stat("/tmp/plugin/snap-plugin-collector-psutil")
		So(fileInfo, ShouldNotBeNil)
		So(err, ShouldBeNil)

	})

	Convey("Should load plugin", t, func() {

		config := ConfigAPI{}
		config.plugins = "collector-psutil"
		config.snapLocation = "/tmp"
		config.SnapPort = "8181"
		execCommand = fakeExecCommand
		err := loadPlugins(config)
		So(err, ShouldBeNil)

	})
}
