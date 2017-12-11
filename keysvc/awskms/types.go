package awskms

import (
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

type AWSKMSService struct {
	client kmsiface.KMSAPI
}
