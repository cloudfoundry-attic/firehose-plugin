package firehose

import (
	"crypto/tls"
	"strconv"
	"sync"

	"fmt"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
)

type Client struct {
	dopplerEndpoint string
	authToken       string
	options         *ClientOptions
	ui              terminal.UI
}

type ClientOptions struct {
	Debug          bool
	NoFilter       bool
	Filter         string
	SubscriptionID string
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
	var err error
	dopplerConnection := consumer.New(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
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
		filter, err = c.promptFilterType()
		if err != nil {
			c.ui.Warn(err.Error())
			return
		}
	}

	subscriptionID := c.options.SubscriptionID
	if len(subscriptionID) == 0 {
		subscriptionID = "FirehosePlugin"
	}

	output, errors := dopplerConnection.FirehoseWithoutReconnect(subscriptionID, c.authToken)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range errors {
			c.ui.Warn(err.Error())
			return
		}
	}()

	defer dopplerConnection.Close()

	c.ui.Say("Starting the nozzle")
	c.ui.Say("Hit Ctrl+c to exit")

	for envelope := range output {
		if filter == "" || filter == strconv.Itoa((int)(envelope.GetEventType())) {
			c.ui.Say("%v \n", envelope)
		}
	}
	wg.Wait()
}

func (c *Client) promptFilterType() (string, error) {

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

	if filter == "" {
		return "", nil
	}

	filterInt, err := strconv.Atoi(filter)
	if err != nil {
		return "", fmt.Errorf("Invalid filter choice %s. Enter an index from 2-9", filter)
	}

	_, ok := events.Envelope_EventType_name[int32(filterInt)]
	if !ok {
		return "", fmt.Errorf("Invalid filter choice %d", filterInt)
	}

	return filter, nil
}

type ConsoleDebugPrinter struct {
	ui terminal.UI
}

func (p ConsoleDebugPrinter) Print(title, dump string) {
	p.ui.Say(title)
	p.ui.Say(dump)
}
