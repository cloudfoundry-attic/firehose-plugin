# nozzle-plugin

[![Build Status](https://travis-ci.org/pivotal-cf-experimental/nozzle-plugin.svg?branch=master)](https://travis-ci.org/pivotal-cf-experimental/nozzle-plugin)

## Installation

```bash
 $ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
 $ cf install-plugin "Firehose Plugin" -r CF-Community

```

## Usage

```bash
cf nozzle --debug (optional)
```

To display all types of messages from the firehose

```bash
cf nozzle --no-filter
```

This only works if logged in as admin

## Uninstall

```bash
cf uninstall FirehosePlugin
```
