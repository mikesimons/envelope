package awskms

import (
	"fmt"
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
	mockDecrypt         func(input *kms.DecryptInput) (*kms.DecryptOutput, error)
}

func (m *mockKMSClient) GenerateDataKey(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
	return m.mockGenerateDataKey(input)
}

func (m *mockKMSClient) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return m.mockDecrypt(input)
}

var _ = Describe("AWSKMS", func() {
	Describe("New", func() {
		It("should return AWS KMS key provider", func() {
			dkp, err := New()
			Expect(err).To(BeNil())
			Expect(dkp).To(BeAssignableToTypeOf(&AWSKMSService{}))
		})
	})

	Describe("GenerateDataKey", func() {
		Context("with valid input", func() {
			It("should create a data key with encryption context", func() {
				context := map[string]string{
					"key1": "val",
					"key2": "val",
				}

				client := &mockKMSClient{
					mockGenerateDataKey: func(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
						Expect(*input.KeyId).To(Equal("testkey"))
						Expect(*input.EncryptionContext["key1"]).To(Equal("val"))
						Expect(*input.EncryptionContext["key2"]).To(Equal("val"))
						return &kms.GenerateDataKeyOutput{CiphertextBlob: []byte("test")}, nil
					},
				}

				dkp := &AWSKMSService{client: client}
				key, err := dkp.GenerateDatakey("testkey", context)

				Expect(err).To(BeNil())
				Expect(key).To(Equal([]byte("test")))
			})
		})

		PIt("should set encryption context in generated key struct")

		Context("with invalid input", func() {
			It("should return error if generateDataKey fails", func() {
				client := &mockKMSClient{
					mockGenerateDataKey: func(input *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
						return &kms.GenerateDataKeyOutput{}, fmt.Errorf("ERROR")
					},
				}

				dkp := &AWSKMSService{client: client}
				_, err := dkp.GenerateDatakey("testkey", nil)

				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("DecryptDatakey", func() {
		It("should set decrypted data key as plaintext field of key", func() {
			ciphertext := []byte("test")
			var plaintext []byte

			client := &mockKMSClient{
				mockDecrypt: func(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
					Expect(input.CiphertextBlob).To(Equal(ciphertext))
					return &kms.DecryptOutput{Plaintext: []byte("abcdefghijklmnopqrstuvwxyzabcdef")}, nil
				},
			}

			dkp := &AWSKMSService{client: client}
			err := dkp.DecryptDatakey(&ciphertext, &plaintext, nil)
			Expect(err).To(BeNil())
			Expect(plaintext).To(Equal([]byte("abcdefghijklmnopqrstuvwxyzabcdef")))
		})

		It("should not decrypt a key twice", func() {
			count := 0
			var ciphertext []byte
			var plaintext []byte

			client := &mockKMSClient{
				mockDecrypt: func(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
					count++
					return &kms.DecryptOutput{Plaintext: []byte("decrypted")}, nil
				},
			}

			dkp := &AWSKMSService{client: client}
			dkp.DecryptDatakey(&ciphertext, &plaintext, nil)
			dkp.DecryptDatakey(&ciphertext, &plaintext, nil)

			Expect(count).To(Equal(1))
		})
	})
})
