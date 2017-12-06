package keysvc

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockGenerateDatakeyResponse struct {
	Output kms.GenerateDataKeyOutput
	Error  error
}

type mockKMSClient struct {
	kmsiface.KMSAPI
	mockGenerateDataKey func(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error)
}

func (m *mockKMSClient) GenerateDataKey(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
	return m.mockGenerateDataKey(input)
}

var _ = Describe("AWSKMS", func() {
	Describe("New", func() {
		It("should return AWS KMS key provider", func() {
			dkp, err := NewAWSKMS()
			Expect(err).To(BeNil())
			Expect(dkp).To(BeAssignableToTypeOf(&AWSKMSService{}))
		})
	})

	Describe("GenerateDataKey", func() {
		Context("with valid input", func() {
			It("should create a data key", func() {
				client := &mockKMSClient{
					mockGenerateDataKey: func(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
						Expect(*input.KeyId).To(Equal("testkey"))
						return &kms.GenerateDataKeyOutput{CiphertextBlob: []byte("test")}, nil
					},
				}

				dkp := &AWSKMSService{client: client}
				key, err := dkp.GenerateDatakey("testkey")

				Expect(err).To(BeNil())
				Expect(key).To(Equal([]byte("test")))
			})

			It("should create a data key with encryption context", func() {
				context := map[string]*string{
					"key1": aws.String("val"),
					"key2": aws.String("val"),
				}

				client := &mockKMSClient{
					mockGenerateDataKey: func(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
						Expect(*input.KeyId).To(Equal("testkey"))
						Expect(input.EncryptionContext).To(Equal(context))
						return &kms.GenerateDataKeyOutput{CiphertextBlob: []byte("test")}, nil
					},
				}

				dkp := &AWSKMSService{client: client}
				key, err := dkp.GenerateDatakeyWithContext("testkey", context)

				Expect(err).To(BeNil())
				Expect(key).To(Equal([]byte("test")))
			})
		})

		Context("with invalid input", func() {
			It("should return error if generateDataKey fails", func() {
				client := &mockKMSClient{
					mockGenerateDataKey: func(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
						return &kms.GenerateDataKeyOutput{}, fmt.Errorf("ERROR")
					},
				}

				dkp := &AWSKMSService{client: client}
				_, err := dkp.GenerateDatakey("testkey")

				Expect(err).ToNot(BeNil())
			})
		})
	})
})
