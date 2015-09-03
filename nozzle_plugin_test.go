package main_test

import (
	"github.com/cloudfoundry/cli/plugin/fakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	. "github.com/jtuchscherer/nozzle-plugin"
	"github.com/jtuchscherer/nozzle-plugin/testhelpers"

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
					nozzlerCmd.Run(fakeCliConnection, []string{"nozzle"})
				})
				outputChan <- output
			}()

			var output []string
			Eventually(outputChan, 2).Should(Receive(&output))

			Expect(output[0]).To(Equal("Starting the nozzle"))
			Expect(output[1]).To(Equal("Hit Cmd+c to exit"))
			Expect(output[2]).To(ContainSubstring("websocket: close 1000"))
			Expect(output[3]).To(ContainSubstring("Log Message"))
		}, 20)
	})

})
