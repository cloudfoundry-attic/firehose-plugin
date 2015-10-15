# nozzle-plugin

[![Build Status](https://travis-ci.org/cloudfoundry/firehose-plugin.svg?branch=master)](https://travis-ci.org/cloudfoundry/firehose-plugin)

## Installation

```bash
 cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
 cf install-plugin "Firehose Plugin" -r CF-Community

```

## Usage
User must be logged in as admin

### Options
```
NAME:
   nozzle - Displays messages from the firehose

USAGE:
   cf nozzle

OPTIONS:
   -debug                 -d, enable debugging
   -filter                -f, specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop
   -no-filter             -n, no firehose filter. Display all messages
   -subscription-id       -s, specify subscription id for distributing firehose output between clients
```

### With Interactive Prompt
```bash
cf nozzle
```

### Without Interactive Prompt
Error message will be displayed for unrecognized filter type

 ```bash
 # For debug
 cf nozzle --debug
 
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
#### Subscription ID
In order to distribute the firehose data evenly among multiple CLI sessions, the user must specify
the same subscription ID to each of the client connections.

 ```bash
 cf nozzle --no-filter --subscription-id myFirehose
 ```

## Uninstall

```bash
cf uninstall FirehosePlugin
```

## Testing

Run tests
```bash
./scripts/test.sh
```

If you want to install the plugin locally and test it manually
```bash
./scripts/install.sh
```
