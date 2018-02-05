package sekrits_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSekrits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sekrits Suite")
}
