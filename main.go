package main

import (
	"os"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/firehose-plugin/firehose"
	"github.com/simonleung8/flags"
)

type NozzlerCmd struct {
	ui terminal.UI
}

func (c *NozzlerCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "FirehosePlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 8,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 3,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "nozzle",
				HelpText: "Displays messages from the firehose",
				UsageDetails: plugin.Usage{
					Usage: "cf nozzle",
					Options: map[string]string{
						"debug":           "-d, enable debugging",
						"no-filter":       "-n, no firehose filter. Display all messages",
						"filter":          "-f, specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop",
						"subscription-id": "-s, specify subscription id for distributing firehose output between clients",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(NozzlerCmd))
}

func (c *NozzlerCmd) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != "nozzle" {
		return
	}

	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	dopplerEndpoint, err := cliConnection.DopplerEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	authToken, err := cliConnection.AccessToken()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	options := c.buildClientOptions(args)

	client := firehose.NewClient(authToken, dopplerEndpoint, options, c.ui)
	client.Start()
}

func (c *NozzlerCmd) buildClientOptions(args []string) *firehose.ClientOptions {
	var debug bool
	var noFilter bool
	var filter string
	var subscriptionId string

	fc := flags.New()
	fc.NewBoolFlag("debug", "d", "used for debugging")
	fc.NewBoolFlag("no-filter", "n", "no firehose filter. Display all messages")
	fc.NewStringFlag("filter", "f", "specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop")
	fc.NewStringFlag("subscription-id", "s", "specify subscription id for distributing firehose output between clients")
	err := fc.Parse(args[1:]...)

	if err != nil {
		c.ui.Failed(err.Error())
	}
	if fc.IsSet("debug") {
		debug = fc.Bool("debug")
	}
	if fc.IsSet("no-filter") {
		noFilter = fc.Bool("no-filter")
	}
	if fc.IsSet("filter") {
		filter = fc.String("filter")
	}
	if fc.IsSet("subscription-id") {
		subscriptionId = fc.String("subscription-id")
	}

	return &firehose.ClientOptions{
		Debug:          debug,
		NoFilter:       noFilter,
		Filter:         filter,
		SubscriptionID: subscriptionId,
	}
}
