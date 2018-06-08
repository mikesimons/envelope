package awskms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"
)

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
