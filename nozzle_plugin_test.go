package main_test

import (
	"strings"

	"github.com/cloudfoundry/cli/plugin/fakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	. "github.com/cloudfoundry/firehose-plugin"
	"github.com/cloudfoundry/firehose-plugin/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	ACCESS_TOKEN = "access_token"
)

var _ = Describe("NozzlePlugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *fakes.FakeCliConnection
		var nozzlerCmd *NozzlerCmd
		var fakeFirehose *testhelpers.FakeFirehose

		BeforeEach(func() {
			fakeFirehose = testhelpers.NewFakeFirehose(ACCESS_TOKEN)
			fakeFirehose.SendLog("Log Message")
			fakeFirehose.Start()

			fakeCliConnection = &fakes.FakeCliConnection{}
			fakeCliConnection.AccessTokenReturns(ACCESS_TOKEN, nil)
			fakeCliConnection.DopplerEndpointReturns(fakeFirehose.URL(), nil)
			nozzlerCmd = &NozzlerCmd{}
		})

		AfterEach(func() {
			fakeFirehose.Close()
		})

		It("works", func(done Done) {
			defer close(done)
			outputChan := make(chan []string)
			go func() {
				output := io_helpers.CaptureOutput(func() {
					nozzlerCmd.Run(fakeCliConnection, []string{"nozzle", "--debug"})
				})
				outputChan <- output
			}()

			var output []string
			Eventually(outputChan, 2).Should(Receive(&output))
			outputString := strings.Join(output, "|")

			Expect(outputString).To(ContainSubstring("What type of firehose messages do you want to see?"))

			Expect(outputString).To(ContainSubstring("Starting the nozzle"))
			Expect(outputString).To(ContainSubstring("Hit Cmd+c to exit"))
			Expect(outputString).To(ContainSubstring("websocket: close 1000"))
			Expect(outputString).To(ContainSubstring("Log Message"))
			Expect(outputString).To(ContainSubstring("WEBSOCKET REQUEST"))
			Expect(outputString).To(ContainSubstring("WEBSOCKET RESPONSE"))
		}, 3)
	})

})
