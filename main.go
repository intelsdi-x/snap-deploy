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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"time"

	"github.com/intelsdi-x/snap-deploy/runner"
	"github.com/urfave/cli"
)

var (
	snapS3URL       = "snap.ci.snap-telemetry.io"
	snapLocation    = ""
	snapteldName    = "snapteld"
	osName          = runtime.GOOS
	taskManifestLoc = "/tmp/task.yml"
	hostname        = "localhost"
	snapteldURL     = fmt.Sprintf("http://%s/snap/latest_build/%s/x86_64/snapteld", snapS3URL, osName)
	snaptelURL      = fmt.Sprintf("http://%s/snap/latest_build/%s/x86_64/snaptel", snapS3URL, osName)

	snapteldDirPath  = filepath.Join(snapLocation, "bin")
	snapteldFilePath = filepath.Join(snapteldDirPath, snapteldName)
	pluginsDirPath   = filepath.Join(snapLocation, "plugin")
	dirLocations     = []string{"bin", "plugin"}
	execCommand      = exec.Command
	removeCommand    = os.RemoveAll
)

func pluginURL(pluginName string) string {

	return fmt.Sprintf("http://snap.ci.snap-telemetry.io/plugins/snap-plugin-%s/latest/%s/x86_64/snap-plugin-%s", pluginName, osName, pluginName)
}

func downloadFromURL(url string, location string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	log.Println("Downloading", url, "to", fileName)

	fullPath := filepath.Join(location, fileName)
	output, err := os.Create(fullPath)
	if err != nil {
		log.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}

	log.Println(n, "bytes downloaded.")
	err = os.Chmod(fullPath, 0755)
	if err != nil {
		log.Println(err)
	}

}

func runSnapd(snapLocation string, port string) int {
	cmd := execCommand("nice", "-10", filepath.Join(snapLocation, "bin", snapteldName), "-t", "0", "-p", port, "-o", "/var/log/")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 7)
	return cmd.Process.Pid

}

func killSnapd() {
	cmd, err := execCommand("killall", snapteldName).CombinedOutput()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(cmd)
	time.Sleep(time.Second * 2)

}

func createDirectories(config ConfigAPI) error {
	err := removeCommand(config.snapLocation)
	if err != nil {

		log.Fatalf("Error removing directory %s\n", err)
	}
	for _, dir := range dirLocations {
		location := filepath.Join(config.snapLocation, dir)
		log.Println("Creating dir: ", location)
		err := os.MkdirAll(location, 0777)
		if err != nil {
			log.Fatalf("Error creating directory %s\n", location)
		}
	}

	return nil
}

func main() {
	app := cli.NewApp()
	config := ConfigAPI{}

	app.Email = "marcin.spoczynski@intel.com"
	app.Usage = "snap-deploy - Snap controller"
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Println("Please use --help switch")
		}
		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dbhost, o",
			Value:       "localhost",
			EnvVar:      "DB_HOST",
			Destination: &config.DbHost,
		},
		cli.StringFlag{
			Name:        "dbdatabase, d",
			Value:       "snap",
			EnvVar:      "DB_NAME",
			Destination: &config.DbDatabase,
		},
		cli.StringFlag{
			Name:        "dbuser, u",
			Value:       "admin",
			EnvVar:      "DB_USER",
			Destination: &config.DbUser,
		},
		cli.StringFlag{
			Name:        "dbpassword, p",
			Value:       "admin",
			EnvVar:      "DB_PASS",
			Destination: &config.DbPassword,
		},
		cli.StringFlag{
			Name:        "interval, i",
			Value:       "1s",
			EnvVar:      "INTERVAL",
			Destination: &config.Interval,
		},
		cli.StringFlag{
			Name:   "tags, t",
			EnvVar: "TAGS",

			Destination: &config.Tags,
		},
		cli.StringFlag{
			Name:        "metrics, m",
			EnvVar:      "METRICS",
			Destination: &config.Metrics,
		},
		cli.StringFlag{
			Name:        "port, sp",
			Value:       "8181",
			EnvVar:      "PORT",
			Destination: &config.SnapPort,
		},
		cli.StringFlag{
			Name:        "directory, e",
			EnvVar:      "DIRECTORY",
			Usage:       "Directory where snap deployed",
			Destination: &config.snapLocation,
		},
		cli.StringFlag{
			Name:        "plugins, r",
			EnvVar:      "PLUGINS",
			Usage:       "Plugins list",
			Destination: &config.plugins,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "deploy",
			Action: func(c *cli.Context) error {
				deploy(config)
				return nil
			},
		},
		{
			Name:    "redeploy",
			Aliases: []string{"rd"},
			Usage:   "redeploy",
			Action: func(c *cli.Context) error {
				killSnapd()
				deploy(config)
				return nil
			},
		},
		{
			Name:    "download",
			Aliases: []string{"f"},
			Usage:   "download",
			Action: func(c *cli.Context) error {
				download(config)
				return nil
			},
		},
		{
			Name:    "kill",
			Aliases: []string{"k"},
			Usage:   "kill",
			Action: func(c *cli.Context) error {
				killSnapd()
				return nil
			},
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start",
			Action: func(c *cli.Context) error {
				runSnapd(config.snapLocation, config.SnapPort)
				return nil
			},
		},
		{
			Name:    "generate_task",
			Aliases: []string{"f"},
			Usage:   "generate_task",
			Action: func(c *cli.Context) error {
				task, err := generateTask(config)
				if err != nil {
					log.Fatalf("Error generating manifest file: %s\n", err.Error())
					panic(err)

				}
				err = ioutil.WriteFile(taskManifestLoc, task, 0644)
				if err != nil {
					log.Fatalf("Error writting manifest file: %s\n", err.Error())
					panic(err)

				}
				log.Fatalf("Manifest file succesfully generated %s", taskManifestLoc)
				return nil
			},
		},
	}
	app.Run(os.Args)

}

func deploy(config ConfigAPI) {
	download(config)
	pid := runSnapd(config.snapLocation, config.SnapPort)
	log.Printf("Snap is running %d\n", pid)
	loadPlugins(config)

	task, err := createTaskCli(config)
	if err != nil {
		log.Fatalf("%s\n", err)

	}
	log.Println("Task created with ID: ", task)

}

func download(config ConfigAPI) {
	err := createDirectories(config)
	if err != nil {
		panic(err)
	}
	downloadFromURL(snapteldURL, filepath.Join(config.snapLocation, "bin"))
	downloadFromURL(snaptelURL, filepath.Join(config.snapLocation, "bin"))

	pluginList := strings.Split(config.plugins, ",")
	for _, plugin := range pluginList {
		downloadFromURL(pluginURL(plugin), filepath.Join(config.snapLocation, "plugin"))
	}

}

func unpackTags(config ConfigAPI) map[string]string {
	var tagMap map[string]string
	tagMap = make(map[string]string)
	pairs := strings.Split(config.Tags, ",")
	for _, p := range pairs {
		item := strings.Split(p, ":")
		if len(item) == 2 {
			tagMap[item[0]] = item[1]
		}
	}
	return tagMap
}

func generateTask(config ConfigAPI) ([]byte, error) {
	schedule := Schedule{
		Type:     "simple",
		Interval: config.Interval}
	publishConfig := PublishConfig{
		Host:     config.DbHost,
		User:     config.DbUser,
		Password: config.DbPassword,
		Database: config.DbDatabase,
		Port:     8086,
	}
	publish := Publish{
		{PluginName: "influxdb", PublishConfig: publishConfig},
	}
	libvirt := Libvirt{
		Nova: true,
	}

	tags := Tags{
		Tags: unpackTags(config),
	}
	configMetric := MetricConfig{
		Libvirt: libvirt,
	}

	collect := Collect{
		Publish:      publish,
		Tags:         tags,
		Metrics:      Metric{},
		MetricConfig: configMetric,
	}

	t := &Task{
		Version:  1,
		Schedule: schedule,
		Workflow: Workflow{
			Collect: collect,
		},
	}
	b, err := json.Marshal(t)
	if err != nil {
		return []byte{}, err

	}
	return createMetricList(config.Metrics, b), nil

}

func createMetricList(metricList string, taskManifest []byte) []byte {
	var metrics []string
	metrics = strings.Split(metricList, ",")
	var buffer bytes.Buffer
	metricLen := len(metrics)
	for i, metric := range metrics {
		metric = strings.Trim(metric, "")
		buffer.WriteString(fmt.Sprintf("\"%s\":{}", metric))
		if i < metricLen-1 {
			buffer.WriteString(",")
		}
	}
	sortedMetrics := strings.Replace(
		string(taskManifest),
		"\"Name\":null", string(buffer.String()),
		1)
	return []byte(sortedMetrics)

}

//Tags info for snap task
type Tags struct {
	Tags map[string]string `json:"/intel"`
}

//ConfigAPI config for Snap Deploy
type ConfigAPI struct {
	DbHost       string
	DbDatabase   string
	DbUser       string
	DbPassword   string
	Tags         string
	Interval     string
	Metrics      string
	Name         string
	SnapCtl      string
	SnapPort     string
	snapLocation string
	plugins      string
}

//Schedule info for snap task
type Schedule struct {
	Type     string `json:"type"`
	Interval string `json:"interval"`
}

//Publish info for snap task
type Publish []struct {
	PluginName    string        `json:"plugin_name"`
	PublishConfig PublishConfig `json:"config"`
}

//PublishConfig info for snap task
type PublishConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

//Process info for snap task
type Process []struct {
	PluginName string      `json:"plugin_name"`
	TagsConfig TagsConfig  `json:"config"`
	ProcessC   interface{} `json:"process"`
	Publish    Publish     `json:"publish"`
}

//TagsConfig info for snap task
type TagsConfig struct {
	Tags string `json:"tags"`
}

//Collect info for snap task
type Collect struct {
	Metrics      Metric       `json:"metrics"`
	Tags         Tags         `json:"tags"`
	Publish      Publish      `json:"publish"`
	MetricConfig MetricConfig `json:"config"`
}

//Task info for snap task
type Task struct {
	Version  int      `json:"version"`
	Schedule Schedule `json:"schedule"`
	Workflow Workflow `json:"workflow"`
}

//Workflow info for snap task
type Workflow struct {
	Collect Collect `json:"collect"`
}

//Metric info for snap task
type Metric struct {
	Name []interface{}
}

//MetricConfig info for snap task
type MetricConfig struct {
	Libvirt Libvirt `json:"/intel/libvirt"`
}

//Libvirt info for snap task
type Libvirt struct {
	Nova bool `json:"nova"`
}

func getSnapURL(config ConfigAPI) string {
	return fmt.Sprintf("http://%s:%s", hostname, config.SnapPort)
}

func loadPlugins(config ConfigAPI) error {

	pluginList := strings.Split(config.plugins, ",")
	snapURL := getSnapURL(config)
	log.Println("Loading plugins")
	for _, plugin := range pluginList {
		log.Println("Loading plugin", fmt.Sprintf("snap-plugin-%s", plugin))
		cmd0, err := execCommand(filepath.Join(config.snapLocation, "bin", "snaptel"), "-u", snapURL, "plugin", "load", filepath.Join(config.snapLocation, "plugin", fmt.Sprintf("snap-plugin-%s", plugin))).CombinedOutput()

		log.Println(string(cmd0))
		if err != nil {
			log.Println(string(err.Error()))
		}
	}
	return nil
}

func createTaskCli(config ConfigAPI) (string, error) {
	fileName := "/tmp/task.yml"
	taskManifest, _ := generateTask(config)
	err := ioutil.WriteFile(fileName, taskManifest, 0644)
	if err != nil {
		log.Fatalf("%s", err)
	}
	url := getSnapURL(config)
	cmdRunner := runner.CmdRunner{}
	reader, err := cmdRunner.Run(filepath.Join(config.snapLocation, "bin", "snaptel"), []string{"-u", url, "task", "create", "-t", fileName})
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ID") {
			fmt.Println(scanner.Text())
			taskID := strings.Split(scanner.Text(), " ")
			return taskID[1], nil

		}

	}

	return "", fmt.Errorf("Can't find task ID")

}
