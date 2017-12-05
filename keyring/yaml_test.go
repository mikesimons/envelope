package keyring_test

import (
	"github.com/mikesimons/sekrits/keyring"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
)

func populateTestFile(realFile string, testFile string, fs afero.Fs) {
	test, _ := os.Open(realFile)
	bytes, _ := ioutil.ReadAll(test)
	afero.WriteFile(fs, testFile, bytes, 0777)
	test.Close()
}

var _ = Describe("YAMLKeyring", func() {
	Describe("LoadYAML", func() {
		Context("with an existing keyring", func() {
			It("should load keys", func() {
				keyring.Fs = afero.NewMemMapFs()
				populateTestFile("testdata/test.yaml", "test.yaml", keyring.Fs)
				kr, err := keyring.LoadYAML("test.yaml")

				keys := kr.GetKeys()

				Expect(err).To(BeNil())
				Expect(keys[0].Name).To(Equal("first"))
			})
		})

		Context("with a new keyring", func() {
			It("should return an empty keyring", func() {
				keyring.Fs = afero.NewMemMapFs()
				kr, err := keyring.LoadYAML("test.yaml")

				keys := kr.GetKeys()

				Expect(err).To(BeNil())
				Expect(keys).To(BeEmpty())
			})
		})

		Context("with invalid input", func() {
			PIt("should raise error if file can't be opened")
			PIt("should raise error if file can't be read")
			PIt("should raise error if file is invalid YAML")
		})
	})

	Describe("AddKey", func() {
		Context("with an existing keyring", func() {
			It("should add a key", func() {
				keyring.Fs = afero.NewMemMapFs()
				populateTestFile("testdata/test.yaml", "test.yaml", keyring.Fs)
				kr, _ := keyring.LoadYAML("test.yaml")

				kr.AddKey(&keyring.Key{Name: "second"})

				keys := kr.GetKeys()
				contents, _ := afero.ReadFile(keyring.Fs, "test.yaml")

				Expect(keys[1].Name).To(Equal("second"))
				Expect(string(contents)).To(ContainSubstring("name: second"))
			})
		})

		Context("with a new keyring", func() {
			It("should add a key", func() {
				keyring.Fs = afero.NewMemMapFs()
				kr, _ := keyring.LoadYAML("test.yaml")

				kr.AddKey(&keyring.Key{Name: "second"})

				keys := kr.GetKeys()
				contents, _ := afero.ReadFile(keyring.Fs, "test.yaml")

				Expect(keys[0].Name).To(Equal("second"))
				Expect(string(contents)).To(ContainSubstring("name: second"))
			})
		})
	})

	Describe("GetKey", func() {
		It("should return the key given a name", func() {
			keyring.Fs = afero.NewMemMapFs()
			populateTestFile("testdata/test.yaml", "test.yaml", keyring.Fs)
			kr, _ := keyring.LoadYAML("test.yaml")

			key, ok := kr.GetKey("first")

			Expect(ok).To(BeTrue())
			Expect(key.Name).To(Equal("first"))
		})

		It("should return the key given an id", func() {
			keyring.Fs = afero.NewMemMapFs()
			populateTestFile("testdata/test.yaml", "test.yaml", keyring.Fs)
			kr, _ := keyring.LoadYAML("test.yaml")

			key, ok := kr.GetKey("first")

			Expect(ok).To(BeTrue())
			Expect(key.Name).To(Equal("first"))
		})

		PIt("should return zero value key and false if key could not be found")
	})

	Describe("GetKeys", func() {
		PIt("should return a list of all keys")
	})
})
