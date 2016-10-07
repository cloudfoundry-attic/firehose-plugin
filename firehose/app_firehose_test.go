package firehose_test

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace/tracefakes"
	"github.com/cloudfoundry/firehose-plugin/firehose"
	"github.com/cloudfoundry/firehose-plugin/firehose/fakes"
	"github.com/cloudfoundry/firehose-plugin/testhelpers"
	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Firehose", func() {
	var (
		ui terminal.UI

		printer      *fakes.FakePrinter
		tracePrinter *tracefakes.FakePrinter

		stdin   *syncedBuffer
		stdout  *syncedBuffer
		options *firehose.ClientOptions
	)

	BeforeEach(func() {
		stdin = &syncedBuffer{}
		stdout = &syncedBuffer{}

		printer = new(fakes.FakePrinter)
		printer.PrintfStub = func(format string, a ...interface{}) (n int, err error) {
			return fmt.Fprintf(stdout, format, a...)
		}
		tracePrinter = new(tracefakes.FakePrinter)

		ui = terminal.NewUI(stdin, stdout, printer, tracePrinter)
		options = &firehose.ClientOptions{AppGUID: "spring-music", Debug: false, NoFilter: true}
	})

	Describe("Start", func() {
		Context("when the connection to doppler cannot be established", func() {
			It("shows a meaningful error", func() {
				client := firehose.NewClient("invalidToken", "badEndpoint", options, ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("Error dialing trafficcontroller server"))
			})
		})
		Context("when the connection to doppler works", func() {
			var fakeFirehose *testhelpers.FakeFirehose
			BeforeEach(func() {
				fakeFirehose = testhelpers.NewFakeFirehoseInAppMode("ACCESS_TOKEN", "spring-music")
				fakeFirehose.SendEvent(events.Envelope_LogMessage, "This is a very special test message")
				fakeFirehose.SendEvent(events.Envelope_ValueMetric, "valuemetric")
				fakeFirehose.SendEvent(events.Envelope_CounterEvent, "counterevent")
				fakeFirehose.SendEvent(events.Envelope_ContainerMetric, "containermetric")
				fakeFirehose.SendEvent(events.Envelope_Error, "this is an error")
				fakeFirehose.SendEvent(events.Envelope_HttpStartStop, "startstop request")
				fakeFirehose.Start()
			})
			It("prints out debug information if demanded", func() {
				options.Debug = true
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("WEBSOCKET REQUEST"))
				Expect(stdout).To(ContainSubstring("WEBSOCKET RESPONSE"))
			})
			It("shows no debug output if not requested", func() {
				options.Debug = false
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
				client.Start()
				Expect(stdout).ToNot(ContainSubstring("WEBSOCKET REQUEST"))
				Expect(stdout).ToNot(ContainSubstring("WEBSOCKET RESPONSE"))
			})
			It("prints out log messages to the terminal", func() {
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("This is a very special test message"))
			})
			Context("in Interactive mode", func() {
				Context("and the user filters by type", func() {
					BeforeEach(func() {
						options.NoFilter = false
					})
					It("does not show log messages when user wants to see ValueMetric", func() {
						stdin.Write([]byte{'6', '\n'})
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).ToNot(ContainSubstring("This is a very special test message"))
					})
					It("shows log messages when the user wants to see log messages", func() {
						stdin.Write([]byte{'5', '\n'})
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("This is a very special test message"))
					})
					It("shows all messages when user hits enter at filter prompt", func() {
						stdin.Write([]byte{'\n'})
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("This is a very special test message"))
						Expect(stdout).To(ContainSubstring("eventType:ValueMetric"))
						Expect(stdout).To(ContainSubstring("eventType:CounterEvent"))
						Expect(stdout).To(ContainSubstring("eventType:ContainerMetric"))
						Expect(stdout).To(ContainSubstring("eventType:Error"))
						Expect(stdout).To(ContainSubstring("eventType:HttpStartStop"))
					})
					It("shows error message when the user enters an invalid filter", func() {
						stdin.Write([]byte{'b', 'l', 'a', '\n'})
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()

						Expect(stdout).To(ContainSubstring("Invalid filter choice bla. Enter an index from 4-9"))
					})
					It("shows error message when the user selects invalid filter index", func() {
						stdin.Write([]byte{'1', '\n'})
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()

						Expect(stdout).To(ContainSubstring("Invalid filter choice 1"))
					})
				})
			})
			Context("in Non-Interactive mode", func() {
				It("errors for un-recognized filter", func() {
					options.NoFilter = false
					options.Filter = "IDontExist"
					stdin.Write([]byte{'1', '\n'})
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()

					Expect(stdout).To(ContainSubstring("Unable to recognize filter IDontExist"))
				})

				It("filters by LogMessage", func() {
					options.Filter = "LogMessage"
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("This is a very special test message"))
				})

				It("filters by ValueMetric", func() {
					options.Filter = "ValueMetric"
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("valueMetric:<name:\"valuemetric\" value:42 unit:\"unit\""))
				})

				It("filters by CounterEvent", func() {
					options.Filter = "CounterEvent"
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("counterEvent:<name:\"counterevent\" delta:42"))
				})

				It("filters by ContainerMetric", func() {
					options.Filter = "ContainerMetric"
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("containerMetric:<applicationId:\"containermetric\" instanceIndex:1 cpuPercentage:1 memoryBytes:1 diskBytes:1"))
				})

				It("filters by Error", func() {
					options.Filter = "Error"
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("error:<source:\"source\" code:404 message:\"this is an error\""))
				})

				It("filters by HttpStartStop", func() {
					options.Filter = "HttpStartStop"
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("httpStartStop:<startTimestamp:1234 stopTimestamp:5555 "))
					Expect(stdout).To(ContainSubstring("userAgent:\"test\""))
					Expect(stdout).To(ContainSubstring("uri:\"http://startstop.example.com\""))
				})

				It("does not filter when NoFilter is true", func() {
					options.NoFilter = true
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(strings.Count(stdout.String(), "eventType:")).To(Equal(6))
				})
			})
		})

	})
})
