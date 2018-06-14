package keysvc_test

import (
	"github.com/mikesimons/envelope/keysvc"
	"github.com/mikesimons/envelope/keysvc/awskms"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("keysvc", func() {
	Describe("GetKeyService", func() {
		It("should return an AWSKMSService given awskms", func() {
			dkp, err := keysvc.GetKeyService("awskms")
			Expect(err).To(BeNil())
			Expect(dkp).To(BeAssignableToTypeOf(&awskms.AWSKMSService{}))
		})

		It("should return an error given an invalid provider", func() {
			_, err := keysvc.GetKeyService("invalid")
			Expect(err).ToNot(BeNil())
		})
	})

	PDescribe("GenerateDatakey", func() {})
})
