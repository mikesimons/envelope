package envelope

import (
	"io"
	"io/ioutil"
)

// Decrypt will decrypt the input as a blob and return the decrypted value as a byte array
func (s *Envelope) Decrypt(input io.Reader) ([]byte, error) {
	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), err
	}

	decrypted, err := s.Keyring.Decrypt(inputBytes)
	if err != nil {
		return []byte(""), err
	}

	return decrypted, nil
}
