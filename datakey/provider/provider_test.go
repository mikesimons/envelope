package provider_test

import (
	"github.com/mikesimons/sekrits/datakey/provider"
	"github.com/mikesimons/sekrits/datakey/provider/awskms"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	Describe("Factory", func() {
		It("should return an awskms.DatakeyProvider given awskms", func() {
			dkp, err := provider.Factory("awskms")
			Expect(err).To(BeNil())
			Expect(dkp).To(BeAssignableToTypeOf(&awskms.DatakeyProvider{}))
		})
	})
})
