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
			return fmt.Sprint(kms.ErrCodeNotFoundException, aerr.Error())
		case kms.ErrCodeDisabledException:
			return fmt.Sprint(kms.ErrCodeDisabledException, aerr.Error())
		case kms.ErrCodeKeyUnavailableException:
			return fmt.Sprint(kms.ErrCodeKeyUnavailableException, aerr.Error())
		case kms.ErrCodeDependencyTimeoutException:
			return fmt.Sprint(kms.ErrCodeDependencyTimeoutException, aerr.Error())
		case kms.ErrCodeInvalidKeyUsageException:
			return fmt.Sprint(kms.ErrCodeInvalidKeyUsageException, aerr.Error())
		case kms.ErrCodeInvalidGrantTokenException:
			return fmt.Sprint(kms.ErrCodeInvalidGrantTokenException, aerr.Error())
		case kms.ErrCodeInternalException:
			return fmt.Sprint(kms.ErrCodeInternalException, aerr.Error())
		case kms.ErrCodeInvalidStateException:
			return fmt.Sprint(kms.ErrCodeInvalidStateException, aerr.Error())
		default:
			return fmt.Sprint(aerr.Error())
		}
	}

	return fmt.Sprintf(err.Error())
}
