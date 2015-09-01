package firehose

import (
	"crypto/tls"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
)

type Client struct {
	dopplerEndpoint string
	authToken       string
	debug           bool
	ui              terminal.UI
}

func NewClient(authToken, doppplerEndpoint string, debug bool, ui terminal.UI) *Client {
	return &Client{
		dopplerEndpoint: doppplerEndpoint,
		authToken:       authToken,
		debug:           debug,
		ui:              ui,
	}

}

func (c *Client) Start() {
	outputChan := make(chan *events.Envelope)
	dopplerConnection := noaa.NewConsumer(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
	if c.debug {
		dopplerConnection.SetDebugPrinter(ConsoleDebugPrinter{ui: c.ui})
	}
	go func() {
		err := dopplerConnection.FirehoseWithoutReconnect("FirehosePlugin", c.authToken, outputChan)
		if err != nil {
			c.ui.Failed(err.Error())
		}
	}()
	defer dopplerConnection.Close()

	c.ui.Say("Starting the nozzle")
	c.ui.Say("Hit Cmd+c to exit")
	c.ui.Say("Hit Cmd+c to exit")
	for msg := range outputChan {
		c.ui.Say("%v \n", msg)
	}
}

type ConsoleDebugPrinter struct {
	ui terminal.UI
}

func (p ConsoleDebugPrinter) Print(title, dump string) {
	p.ui.Say(title)
	p.ui.Say(dump)
}
