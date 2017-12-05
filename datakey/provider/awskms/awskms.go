package awskms

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	//errors "github.com/hashicorp/errwrap"
)

var awsSession = session.Must(session.NewSession())

type GenerateDatakey struct {
}

type DatakeyProvider struct {
	client kmsiface.KMSAPI
}

func New() (*DatakeyProvider, error) {
	kmsClient := kms.New(awsSession, aws.NewConfig())
	return &DatakeyProvider{
		client: kmsClient,
	}, nil
}

func (dkp *DatakeyProvider) DecryptDatakey(key []byte) {

}

func (dkp *DatakeyProvider) GenerateDatakey(key string) ([]byte, error) {
	return dkp.GenerateDatakeyWithContext(key, nil)
}

func (dkp *DatakeyProvider) GenerateDatakeyWithContext(key string, context map[string]*string) ([]byte, error) {
	input := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(key),
		KeySpec: aws.String("AES_256"),
	}

	if context != nil {
		input.EncryptionContext = context
	}

	datakey, err := dkp.client.GenerateDataKey(input)
	if err != nil {

		return nil, fmt.Errorf("Couldn't generate a data key using KMS: %v", awsErrorString(err))
	}

	return datakey.CiphertextBlob, nil
}

func awsErrorString(err error) string {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case kms.ErrCodeNotFoundException:
			return fmt.Sprintf(kms.ErrCodeNotFoundException, aerr.Error())
		case kms.ErrCodeDisabledException:
			return fmt.Sprintf(kms.ErrCodeDisabledException, aerr.Error())
		case kms.ErrCodeKeyUnavailableException:
			return fmt.Sprintf(kms.ErrCodeKeyUnavailableException, aerr.Error())
		case kms.ErrCodeDependencyTimeoutException:
			return fmt.Sprintf(kms.ErrCodeDependencyTimeoutException, aerr.Error())
		case kms.ErrCodeInvalidKeyUsageException:
			return fmt.Sprintf(kms.ErrCodeInvalidKeyUsageException, aerr.Error())
		case kms.ErrCodeInvalidGrantTokenException:
			return fmt.Sprintf(kms.ErrCodeInvalidGrantTokenException, aerr.Error())
		case kms.ErrCodeInternalException:
			return fmt.Sprintf(kms.ErrCodeInternalException, aerr.Error())
		case kms.ErrCodeInvalidStateException:
			return fmt.Sprintf(kms.ErrCodeInvalidStateException, aerr.Error())
		default:
			return fmt.Sprintf(aerr.Error())
		}
	}

	return fmt.Sprintf(err.Error())
}
