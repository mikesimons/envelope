package awskms

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	errors "github.com/hashicorp/errwrap"
)

var awsSession = session.Must(session.NewSession())

func New() (*AWSKMSService, error) {
	kmsClient := kms.New(awsSession, aws.NewConfig())
	return &AWSKMSService{
		client: kmsClient,
	}, nil
}

func (dkp *AWSKMSService) DecryptDatakey(ciphertext *[]byte, plaintext *[]byte) error {
	if len(*plaintext) > 0 {
		return nil
	}

	input := &kms.DecryptInput{
		CiphertextBlob: *ciphertext,
	}
	result, err := dkp.client.Decrypt(input)

	if err != nil {
		return errors.Wrapf("Could not decrypt AWS KMS data key", err)
	}

	ret := make([]byte, len(result.Plaintext))
	copy(ret[:], result.Plaintext[:])
	*plaintext = ret

	return nil
}

func (dkp *AWSKMSService) GenerateDatakey(key string) ([]byte, error) {
	return dkp.GenerateDatakeyWithContext(key, nil)
}

func (dkp *AWSKMSService) GenerateDatakeyWithContext(key string, context map[string]*string) ([]byte, error) {
	input := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(key),
		KeySpec: aws.String("AES_256"),
	}

	if context != nil {
		input.EncryptionContext = context
	}

	datakey, err := dkp.client.GenerateDataKey(input)
	if err != nil {
		return []byte(""), errors.Wrapf(fmt.Sprintf("Couldn't generate a data key using AWS KMS: %v", awsErrorString(err)), err)
	}

	return datakey.CiphertextBlob, nil
}
