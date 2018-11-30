package envelope

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

type EncryptOpts struct {
	Encoder    func([]byte, io.Writer) error
	WithPrefix bool
}

func NopEncoder(in []byte, buf io.Writer) error {
	_, err := io.Copy(buf, bytes.NewReader(in))
	return err
}

func Base64Encoder(in []byte, buf io.Writer) error {
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	_, err := encoder.Write(in)
	encoder.Close()
	return err
}

// EncryptWithOpts will encypt the input as a blob using the key given and the options specified
func (s *Envelope) EncryptWithOpts(alias string, input io.Reader, opts EncryptOpts) ([]byte, error) {
	key, ok := s.Keyring.GetKey(alias)
	if !ok {
		return []byte(""), fmt.Errorf("Couldn't find key with alias or id '%s'", alias)
	}

	data, err := ioutil.ReadAll(input)
	if err != nil {
		return []byte(""), err
	}

	encrypted, err := key.EncryptBytes(data)
	if err != nil {
		return []byte(""), err
	}

	var output []byte
	buffer := bytes.NewBuffer(output)
	if opts.WithPrefix {
		buffer.Write([]byte(s.Prefix))
	}

	err = opts.Encoder(encrypted, buffer)
	if err != nil {
		return []byte(""), err
	}

	return ioutil.ReadAll(buffer)
}

// Encrypt will encrypt the input as a blob using the key given and default options
func (s *Envelope) Encrypt(alias string, input io.Reader) ([]byte, error) {
	return s.EncryptWithOpts(alias, input, EncryptOpts{Encoder: NopEncoder})
}
