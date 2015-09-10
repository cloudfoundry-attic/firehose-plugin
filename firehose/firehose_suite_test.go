package firehose_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFirehose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Firehose Suite")
}
