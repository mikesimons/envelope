package keyring_test

import (
	"github.com/mikesimons/sekrits/keyring"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keyring", func() {
	Describe("Load", func() {
		It("should return a YAMLKeyring when given a yaml file path", func() {
			kr, err := keyring.Load("test.yaml")
			Expect(kr).To(BeAssignableToTypeOf(&keyring.YAMLKeyring{}))
			Expect(err).To(BeNil())
		})
	})
})
