package firehose_test

import (
	"fmt"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/terminal/fakes"
	"github.com/jtuchscherer/nozzle-plugin/firehose"
	"github.com/jtuchscherer/nozzle-plugin/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeStdin struct {
	Input string
}

func (r fakeStdin) Read(b []byte) (n int, err error) {
	b = append(b, []byte(r.Input)...)
	return len(r.Input), nil
}

var _ = Describe("Firehose", func() {
	var printer *fakes.FakePrinter
	var ui terminal.UI
	var stdin fakeStdin
	var stdout string
	BeforeEach(func() {
		printer = new(fakes.FakePrinter)
		stdout = ""
		printer.PrintfStub = func(format string, a ...interface{}) (n int, err error) {
			stdout += fmt.Sprintf(format, a...)
			return len(stdout), nil
		}
		stdin = fakeStdin{}
		ui = terminal.NewUI(stdin, printer)
	})

	Context("Start", func() {
		Context("when the connection to doppler cannot be established", func() {
			It("shows a meaningful error", func() {
				client := firehose.NewClient("invalidToken", "badEndpoint", false, ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("Error dialing traffic controller server"))
			})

		})
		Context("when the connection to doppler works", func() {
			var fakeFirehose *testhelpers.FakeFirehose
			BeforeEach(func() {
				fakeFirehose = testhelpers.NewFakeFirehose("ACCESS_TOKEN")
				fakeFirehose.SendLog("This is a very special test message")
				fakeFirehose.Start()
			})
			It("prints out debug information if demanded", func() {
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), true, ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("WEBSOCKET REQUEST"))
				Expect(stdout).To(ContainSubstring("WEBSOCKET RESPONSE"))
			})
			It("shows no debug output if not requested", func() {
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), false, ui)
				client.Start()
				Expect(stdout).ToNot(ContainSubstring("WEBSOCKET REQUEST"))
				Expect(stdout).ToNot(ContainSubstring("WEBSOCKET RESPONSE"))
			})
			It("prints out log messages to the terminal", func() {
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), false, ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("This is a very special test message"))
			})
		})
	})
})
