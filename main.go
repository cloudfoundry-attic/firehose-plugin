package main

import (
	"os"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin"
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
			Minor: 11,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 17,
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
			{
				Name:     "app-nozzle",
				HelpText: "Displays messages from the firehose for a given app",
				UsageDetails: plugin.Usage{
					Usage: "cf app-nozzle APP_NAME",
					Options: map[string]string{
						"debug":     "-d, enable debugging",
						"no-filter": "-n, no filter. Display all messages",
						"filter":    "-f, specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop",
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
	var options *firehose.ClientOptions

	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	switch args[0] {
	case "nozzle":
		options = c.buildClientOptions(args)
	case "app-nozzle":
		options = c.buildClientOptions(args)
		appModel, err := cliConnection.GetApp(args[1])
		if err != nil {
			c.ui.Warn(err.Error())
			return
		}

		options.AppGUID = appModel.Guid
	default:
		return
	}

	dopplerEndpoint, err := cliConnection.DopplerEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	authToken, err := cliConnection.AccessToken()
	if err != nil {
		c.ui.Failed(err.Error())
	}

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
