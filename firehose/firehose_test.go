package firehose_test

import (
	"fmt"
	"io"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/terminal/fakes"
	"github.com/jtuchscherer/nozzle-plugin/firehose"
	"github.com/jtuchscherer/nozzle-plugin/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeStdin struct {
	Input []byte
	done  bool
}

func (r *fakeStdin) Read(p []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}
	for i, b := range r.Input {
		p[i] = b
	}
	r.done = true
	return len(r.Input), nil
}

var _ = Describe("Firehose", func() {
	var printer *fakes.FakePrinter
	var ui terminal.UI
	var stdin *fakeStdin
	var stdout string

	BeforeEach(func() {
		printer = new(fakes.FakePrinter)
		stdout = ""
		printer.PrintfStub = func(format string, a ...interface{}) (n int, err error) {
			stdout += fmt.Sprintf(format, a...)
			return len(stdout), nil
		}
		stdin = &fakeStdin{[]byte{'\n'}, false}
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
			Context("and the user filters by type", func() {
				It("does not show log messages when user wants to see HttpStart", func() {
					stdin.Input = []byte{'2', '\n'}
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), false, ui)
					client.Start()
					Expect(stdout).ToNot(ContainSubstring("This is a very special test message"))
				})
				It("shows log messages when the user wants to see log messages", func() {
					stdin.Input = []byte{'5', '\n'}
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), false, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("This is a very special test message"))
				})

			})
		})
	})
})
