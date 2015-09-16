package firehose_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	r   *os.File
	w   *os.File
	old *os.File
)

func TestFirehose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Firehose Suite")
}

var _ = BeforeSuite(func() {
	old = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w
})

var _ = AfterSuite(func() {
	w.Close()
	os.Stdout = old
})
