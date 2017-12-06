package datakey_test

import (
	"github.com/mikesimons/sekrits/datakey"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	Describe("Factory", func() {
		It("should return an awskms.DatakeyProvider given awskms", func() {
			dkp, err := datakey.Factory("awskms")
			Expect(err).To(BeNil())
			Expect(dkp).To(BeAssignableToTypeOf(&datakey.AWSKMSProvider{}))
		})
	})
})
