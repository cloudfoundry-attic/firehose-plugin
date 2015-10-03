package firehose

import (
	"crypto/tls"
	"strconv"

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

	filter := c.ui.Ask(`What type of firehose messages do you want to see? Please enter one of the following choices:
  hit 'enter' for all messages
  2 for HttpStart
  3 for HttpStop
  4 for HttpStartStop
  5 for LogMessage
  6 for ValueMetric
  7 for CounterEvent
  8 for Error
  9 for ContainerMetric
`)

	go func() {
		err := dopplerConnection.FirehoseWithoutReconnect("FirehosePlugin", c.authToken, outputChan)
		if err != nil {
			c.ui.Warn(err.Error())
			close(outputChan)
			return
		}
	}()

	defer dopplerConnection.Close()

	c.ui.Say("Starting the nozzle")
	c.ui.Say("Hit Ctrl+c to exit")

	for envelope := range outputChan {
		if filter == "" || filter == strconv.Itoa((int)(envelope.GetEventType())) {
			c.ui.Say("%v \n", envelope)
		}
	}
}

type ConsoleDebugPrinter struct {
	ui terminal.UI
}

func (p ConsoleDebugPrinter) Print(title, dump string) {
	p.ui.Say(title)
	p.ui.Say(dump)
}
