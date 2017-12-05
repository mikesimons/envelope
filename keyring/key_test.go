package keyring_test

import (
	"github.com/mikesimons/sekrits/keyring"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Key", func() {
	Describe("NewKey", func() {
		It("should return a new key with alias, cipher text and id", func() {
			key := keyring.NewKey("alias", []byte("ciphertext"))
			Expect(key.Alias).To(Equal("alias"))
			Expect(key.Ciphertext).To(Equal([]byte("ciphertext")))
			Expect(key.Id).NotTo(BeAssignableToTypeOf(&uuid.UUID{}))
		})
	})

	Describe("Custom YAML", func() {
		Describe("Marshal", func() {
			It("should encode id, alias & ciphertext", func() {
				key := &keyring.Key{
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
	})
})
