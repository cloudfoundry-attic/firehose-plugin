package main_test

import (
	"github.com/cloudfoundry/cli/testhelpers/pluginbuilder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNozzlePlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	pluginbuilder.BuildTestBinary(".", "main")
	RunSpecs(t, "NozzlePlugin Suite")
}
