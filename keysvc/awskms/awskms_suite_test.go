package awskms_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAwskms(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Awskms Suite")
}
