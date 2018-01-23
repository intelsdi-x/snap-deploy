[![Go Report Card](http://goreportcard.com/badge/intelsdi-x/snap-deploy)](http://goreportcard.com/report/intelsdi-x/snap-deploy)

# snap-deploy

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

Linux/MacOS/*BSD system

#### Download binary:

You can get the pre-built binaries for your OS and architecture from the [GitHub Releases](https://github.com/intelsdi-x/snap-deploy/releases) page. Download the plugin from the latest release and load it into `/opt/local/bin` is the default location for user binaries.

#### To build the snap-deploy binary:

Fork https://github.com/intelsdi-x/snap-deploy
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-deploy.git
```

Build the snap-deploy by running make within the cloned repo:
```
$ make
```
This builds the executable file in `./build/`

## Documentation

### Examples

For quick snap-deploy test using deployment, you can go through steps below:

1. Install influxdb on your localhost [InfluxDB website](https://www.influxdata.com/)
2. Create database "snap"
3. Export variables:
```
export DB_HOST="localhost"
export DB_NAME="snap"
export DB_USER="snap"
export DB_PASS="snap"
export INTERVAL="1s"
export TAGS="datacenter:dublin"
export METRICS="/intel"
export PORT="8181"
export DIRECTORY="/opt/snap"
export PLUGINS="collector-psutil"

```

4. Run snap-deploy as root
```
snap-deploy deploy
```

5. Snap should be running on the port 8181 and publishing psutil metrics to the localhost influxdb database named snap. 

### Roadmap
As we launch this software, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-deploy/issues).

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-deploy/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-deploy/pulls).

## Community Support

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com/intelsdi-x/snap), along with this application, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
This is Open Source software released under the Apache 2.0 License. Please see the [LICENSE](LICENSE) file for full license details.

* Author: [Marcin Spoczynski](https://github.com/sandlbn/)

This software has been contributed by MIKELANGELO, a Horizon 2020 project co-funded by the European Union. https://www.mikelangelo-project.eu/
## Thank You
And **thank you!** Your contribution, through code and participation, is incredibly important to us.
