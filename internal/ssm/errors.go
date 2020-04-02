package ssm

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
)

func CheckSSMGetParameters(err error) {
	var errstr string

	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInvalidKeyId:
			errstr = awsErr.Error()
		case awsssm.ErrCodeInternalServerError:
			errstr = awsErr.Error()
		}
		fmt.Printf("SSM Get Parameters Error: %s", errstr)
	}
}

func CheckSSMError(err error) {
	var errstr string

	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInternalServerError:
			errstr = awsErr.Error()
		case awsssm.ErrCodeInvalidKeyId:
			errstr = awsErr.Error()
		case awsssm.ErrCodeParameterNotFound:
			errstr = awsErr.Error()
		case awsssm.ErrCodeParameterVersionNotFound:
			errstr = awsErr.Error()
		}
		fmt.Printf("Error: %s", errstr)
	}
}

func CheckSSMByPath(err error) {
	var errstr string
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInternalServerError:
			errstr = awsErr.Error()
		case awsssm.ErrCodeInvalidFilterKey:
			errstr = awsErr.Error()
		case awsssm.ErrCodeInvalidFilterOption:
			errstr = awsErr.Error()
		case awsssm.ErrCodeInvalidFilterValue:
			errstr = awsErr.Error()
		case awsssm.ErrCodeInvalidKeyId:
			errstr = awsErr.Error()
		case awsssm.ErrCodeInvalidNextToken:
			errstr = awsErr.Error()
		}
		fmt.Printf("SSM By Path Error: %s", errstr)
	}
}
