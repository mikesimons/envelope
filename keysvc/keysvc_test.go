package keysvc_test

import (
	"github.com/mikesimons/sekrits/keysvc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	Describe("Factory", func() {
		It("should return an AWSKMSService given awskms", func() {
			dkp, err := keysvc.Factory("awskms")
			Expect(err).To(BeNil())
			Expect(dkp).To(BeAssignableToTypeOf(&keysvc.AWSKMSService{}))
		})

		It("should return an error given an invalid provider", func() {
			_, err := keysvc.Factory("invalid")
			Expect(err).ToNot(BeNil())
		})
	})
})
