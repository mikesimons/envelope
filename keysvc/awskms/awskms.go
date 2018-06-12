package awskms

import (
	"strings"

	"github.com/ansel1/merry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

var awsSession = session.Must(session.NewSession())

func New() (*AWSKMSService, error) {
	kmsClient := kms.New(awsSession, aws.NewConfig())
	return &AWSKMSService{
		client: kmsClient,
	}, nil
}

func (dkp *AWSKMSService) DecryptDatakey(ciphertext *[]byte, plaintext *[]byte, context map[string]string) error {
	if len(*plaintext) > 0 {
		return nil
	}

	input := &kms.DecryptInput{
		CiphertextBlob:    *ciphertext,
		EncryptionContext: convertContext(context),
	}
	result, err := dkp.client.Decrypt(input)

	if err != nil {
		return merry.Wrap(err).WithUserMessage("Could not decrypt AWS KMS data key")
	}

	ret := make([]byte, len(result.Plaintext))
	copy(ret[:], result.Plaintext[:])
	*plaintext = ret

	return nil
}

func (dkp *AWSKMSService) GenerateDatakey(key string, context map[string]string) ([]byte, error) {
	input := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(key),
		KeySpec: aws.String("AES_256"),
	}

	if context != nil {
		input.EncryptionContext = convertContext(context)
	}

	datakey, err := dkp.client.GenerateDataKey(input)
	if err != nil {
		return []byte(""), merry.Wrap(err).
			WithUserMessage("Couldn't generate a data key using AWS KMS").
			WithMessage(strings.Replace(awsErrorString(err), "\n", "", -1))
	}

	return datakey.CiphertextBlob, nil
}

func convertContext(input map[string]string) map[string]*string {
	values := make([]string, len(input))
	output := make(map[string]*string)
	for k, v := range input {
		values = append(values, v)
		output[k] = &values[len(values)-1]
	}
	return output
}
