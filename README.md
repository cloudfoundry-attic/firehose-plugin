# nozzle-plugin

[![Build Status](https://travis-ci.org/pivotal-cf-experimental/nozzle-plugin.svg?branch=master)](https://travis-ci.org/pivotal-cf-experimental/nozzle-plugin)

## Installation

```bash
 $ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
 $ cf install-plugin "Firehose Plugin" -r CF-Community

```

## Usage

### With Interactive Prompt
```bash
cf nozzle --debug (optional)
```

### Without Interactive Prompt
- For invalid message types, it will prompt for input.
- This only works if logged in as admin
- To filter out messages based on type
 ```bash
 # For all messages
 cf nozzle --no-filter
 
 # For Log Messages
 cf nozzle --filter LogMessage
 
 # For HttpStart
 cf nozzle --filter HttpStart
 
 # For HttpStartStop
 cf nozzle --filter HttpStartStop
 
 # For HttpStop
 cf nozzle --filter HttpStop
 
 # For ValueMetric
 cf nozzle --filter ValueMetric
 
 # For CounterEvent
 cf nozzle --filter CounterEvent
 
 # For ContainerMetric
 cf nozzle --filter ContainerMetric
 
 # For Error
 cf nozzle --filter Error
 ```

## Uninstall

```bash
cf uninstall FirehosePlugin
```
