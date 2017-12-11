package keysvc

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

type TestKeyService struct {
	mockGenerateDatakey func(string) ([]byte, error)
	mockDecryptDatakey  func(*[]byte, *[]byte) error
}

func (t *TestKeyService) GenerateDatakey(master string) ([]byte, error) {
	return t.mockGenerateDatakey(master)
}

func (t *TestKeyService) DecryptDatakey(ciphertext *[]byte, plaintext *[]byte) error {
	return t.mockDecryptDatakey(ciphertext, plaintext)
}

func NewTestKeyService() (KeyServiceProvider, error) {
	return &TestKeyService{
		mockGenerateDatakey: func(master string) ([]byte, error) {
			return []byte(""), nil
		},
		mockDecryptDatakey: func(ciphertext *[]byte, plaintext *[]byte) error {
			*plaintext = make([]byte, 32)
			return nil
		},
	}, nil
}

var _ = Describe("Key", func() {
	Describe("NewKey", func() {
		It("should return a new key with alias, cipher text and id", func() {
			key := NewKey("alias", "test", []byte("ciphertext"))
			Expect(key.Alias).To(Equal("alias"))
			Expect(key.Ciphertext).To(Equal([]byte("ciphertext")))
			Expect(key.Id).NotTo(BeAssignableToTypeOf(&uuid.UUID{}))
		})
	})

	Describe("EncryptBytes / Decrypt", func() {
		It("should encrypt given data with key in a way that can be decrypted only given encrypted blob", func() {
			AddKeyServiceFn("test", NewTestKeyService)

			key := NewKey("alias", "test", []byte("test"))

			encrypted, err := key.EncryptBytes([]byte("hello"))
			Expect(err).To(BeNil())

			decrypted, err := key.Decrypt(encrypted)
			Expect(err).To(BeNil())

			Expect(decrypted).To(Equal([]byte("hello")))
		})
	})

	Describe("Custom YAML", func() {
		Describe("Marshal", func() {
			It("should encode id, alias & ciphertext", func() {
				key := &Key{
					Id:         uuid.NewV4(),
					Alias:      "alias",
					Ciphertext: []byte("test"),
				}

				marshalled, err := yaml.Marshal(key)

				Expect(err).To(BeNil())

				verify := make(map[string]string)
				err = yaml.Unmarshal(marshalled, &verify)

				Expect(err).To(BeNil())
				Expect(verify["id"]).To(Equal(key.Id.String()))
				Expect(verify["alias"]).To(Equal("alias"))
				Expect(verify["key"]).To(Equal("dGVzdA==")) // Base64 encoded "test"
			})
		})

		Describe("Unmarshal", func() {
			It("should decide id, alias & ciphertext", func() {
				inputKey := &Key{
					Id:         uuid.NewV4(),
					Alias:      "alias",
					Ciphertext: []byte("test"),
				}

				marshalled, err := yaml.Marshal(inputKey)

				Expect(err).To(BeNil())

				outputKey := &Key{}
				err = yaml.Unmarshal(marshalled, &outputKey)

				Expect(err).To(BeNil())
				Expect(outputKey).To(Equal(inputKey))
			})
		})
	})
})
