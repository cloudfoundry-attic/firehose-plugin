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
	options         *ClientOptions
	ui              terminal.UI
}

type ClientOptions struct {
	Debug    bool
	NoFilter bool
	Filter   string
}

func NewClient(authToken, doppplerEndpoint string, options *ClientOptions, ui terminal.UI) *Client {
	return &Client{
		dopplerEndpoint: doppplerEndpoint,
		authToken:       authToken,
		options:         options,
		ui:              ui,
	}

}

func (c *Client) Start() {
	outputChan := make(chan *events.Envelope)
	dopplerConnection := noaa.NewConsumer(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
	if c.options.Debug {
		dopplerConnection.SetDebugPrinter(ConsoleDebugPrinter{ui: c.ui})
	}
	filter := ""
	switch {
	case c.options.NoFilter:
		filter = ""
	case c.options.Filter != "":
		envelopeType, ok := events.Envelope_EventType_value[c.options.Filter]
		if !ok {
			c.ui.Warn("Unable to recognize filter %s", c.options.Filter)
			return
		}
		filter = strconv.Itoa(int(envelopeType))

	default:
		c.ui.Say("What type of firehose messages do you want to see?")
		filter = c.promptFilterType()
	}

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

func (c *Client) promptFilterType() string {

	filter := c.ui.Ask(`Please enter one of the following choices:
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

	return filter
}

type ConsoleDebugPrinter struct {
	ui terminal.UI
}

func (p ConsoleDebugPrinter) Print(title, dump string) {
	p.ui.Say(title)
	p.ui.Say(dump)
}
