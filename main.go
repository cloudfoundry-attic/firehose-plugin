package main

import (
	"os"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/flags"
	"github.com/cloudfoundry/cli/flags/flag"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/firehose-plugin/firehose"
)

type NozzlerCmd struct {
	ui terminal.UI
}

func (c *NozzlerCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "FirehosePlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 5,
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
				HelpText: "Command to print out messages from the firehose",
				UsageDetails: plugin.Usage{
					Usage: "cf nozzle",
					Options: map[string]string{
						"debug": "true to enable debugging",
					},
				},
			},
		},
	}
}

func setupFlags() map[string]flags.FlagSet {
	fs := make(map[string]flags.FlagSet)
	fs["debug"] = &cliFlags.BoolFlag{Name: "debug", Usage: "used for debugging"}
	return fs
}

func main() {
	plugin.Start(new(NozzlerCmd))
}

func (c *NozzlerCmd) Run(cliConnection plugin.CliConnection, args []string) {
	var debug bool
	if args[0] != "nozzle" {
		return
	}
	c.ui = terminal.NewUI(os.Stdin, terminal.NewTeePrinter())

	fc := flags.NewFlagContext(setupFlags())
	err := fc.Parse(args[1:]...)
	if err != nil {
		c.ui.Failed(err.Error())
	}
	if fc.IsSet("debug") {
		debug = fc.Bool("debug")
	}

	dopplerEndpoint, err := cliConnection.DopplerEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	authToken, err := cliConnection.AccessToken()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	client := firehose.NewClient(authToken, dopplerEndpoint, debug, c.ui)
	client.Start()
}
