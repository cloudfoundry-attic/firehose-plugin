# nozzle-plugin

[![Build Status](https://travis-ci.org/pivotal-cf-experimental/nozzle-plugin.svg?branch=master)](https://travis-ci.org/pivotal-cf-experimental/nozzle-plugin)

## Installation

```bash
# In your go workspace:

go get github.com/pivotal-cf-experimental/nozzle-plugin
cf install-plugin ./bin/nozzle-plugin

```

## Usage

```bash
cf nozzle --debug (optional)
```
This only works if logged in as admin

## Uninstall

```bash
cf uninstall FirehosePlugin
```
