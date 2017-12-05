package keyring_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKeyring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keyring Suite")
}
