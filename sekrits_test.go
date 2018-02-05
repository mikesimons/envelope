package sekrits_test

import (
	"github.com/mikesimons/sekrits/keyring"
	"github.com/mikesimons/sekrits"

	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"os"
)

var testKeyId = os.Getenv("SEKRITS_KMS_TEST_KEY_ID")
var testKeyAlias = os.Getenv("SEKRITS_KMS_TEST_KEY_ALIAS")

var _ = Describe("Sekrits", func() {
	Describe("AddKey", func() {
		It("should add a key to the keyring", func() {
			keyring.Fs = afero.NewMemMapFs()
			app, _ := sekrits.WithYamlKeyring("test.yaml")
			keyId, err := app.AddKey("alias", fmt.Sprintf("awskms://%s", testKeyId))

			Expect(keyId).ToNot(BeEmpty())
			Expect(err).To(BeNil())

			file, _ := afero.ReadFile(keyring.Fs, "test.yaml")
			var output []map[string]string
			yaml.Unmarshal(file, &output)

			val := output[0]
			Expect(val["alias"]).To(Equal("alias"))
		})

		It("should return an error if the alias clashes with another key", func() {
			keyring.Fs = afero.NewMemMapFs()
			app, _ := sekrits.WithYamlKeyring("test.yaml")

			keyId, err := app.AddKey("alias", fmt.Sprintf("awskms://%s", testKeyId))
			Expect(err).To(BeNil())

			_, err = app.AddKey("alias", fmt.Sprintf("awskms://%s", testKeyId))
			Expect(err).ToNot(BeNil())

			_, err = app.AddKey(keyId, fmt.Sprintf("awskms://%s", testKeyId))
			Expect(err).ToNot(BeNil())
		})

		It("should return an error if the data key could not be generated", func() {
			keyring.Fs = afero.NewMemMapFs()
			app, _ := sekrits.WithYamlKeyring("test.yaml")
			_, err := app.AddKey("alias", "awskms://")
			Expect(err).ToNot(BeNil())
		})

		It("should return an error if an invalid dsn is provided", func() {
			keyring.Fs = afero.NewMemMapFs()
			app, _ := sekrits.WithYamlKeyring("test.yaml")
			_, err := app.AddKey("alias", "://://")
			Expect(err).ToNot(BeNil())
		})

		It("should return an error if an invalid key service is provided", func() {
			keyring.Fs = afero.NewMemMapFs()
			app, _ := sekrits.WithYamlKeyring("test.yaml")
			_, err := app.AddKey("alias", fmt.Sprintf("kms://%s", testKeyId))
			Expect(err).ToNot(BeNil())
		})

		It("should return an error if an existing keyring can't be loaded", func() {
			keyring.Fs = afero.NewMemMapFs()
			afero.WriteFile(keyring.Fs, "test.yaml", []byte("this\nis\nnot\nvalid\nyaml"), 0644)
			_, err := sekrits.WithYamlKeyring("test.yaml")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Encrypt / Decrypt", func() {
		It("should encrypt the given secret in a way that can be decrypted", func() {
			keyring.Fs = afero.NewMemMapFs()
			app, _ := sekrits.WithYamlKeyring("test.yaml")
			_, err := app.AddKey("alias", fmt.Sprintf("awskms://%s", testKeyId))
			Expect(err).To(BeNil())

			encrypted, err := app.Encrypt("alias", bytes.NewReader([]byte("teststring")))
			Expect(err).To(BeNil())

			decrypted, err := app.Decrypt(bytes.NewReader(encrypted))
			Expect(err).To(BeNil())
			Expect(decrypted).To(Equal([]byte("teststring")))
		})

		PIt("should return an error if an invalid key is given")

		PIt("should return an error if the input can not be read")

		PIt("should return an error if encryption fails")

		It("should return an error if an existing keyring can't be loaded", func() {
			keyring.Fs = afero.NewMemMapFs()
			afero.WriteFile(keyring.Fs, "test.yaml", []byte("this\nis\nnot\nvalid\nyaml"), 0644)
			_, err := sekrits.WithYamlKeyring("test.yaml")
			Expect(err).ToNot(BeNil())
		})
	})

	PDescribe("Decrypt structured", func() {
		PIt("should decrypt all encrypted values in a yaml document")
		PIt("should decrypt all encrypted values in a json document")
		PIt("should decrypt all encrypted values in a toml document")
		PIt("should return an error if yaml document can't be parsed")
		PIt("should return an error if json document can't be parsed")
		PIt("should return an error if toml document can't be parsed")
		PIt("should return an error if document type is blob or not recognized")
		PIt("should return an error if a value could not be decrypted")
		PIt("should unset value if a value could not be decrypted and ignore decryption errors is set")
	})
})
