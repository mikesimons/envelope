package keysvc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestProvider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Key Service Suite")
}
