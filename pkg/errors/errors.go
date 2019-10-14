package errors

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func CheckSSMGetParameters(err error) {
	var errstr string

	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ssm.ErrCodeInvalidKeyId:
			errstr = awsErr.Error()
		case ssm.ErrCodeInternalServerError:
			errstr = awsErr.Error()
		}
	}
	fmt.Printf("SSM Get Parameters Error: %s", errstr)
}

func CheckSSMError(err error) {
	var errstr string

	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ssm.ErrCodeInternalServerError:
			errstr = awsErr.Error()
		case ssm.ErrCodeInvalidKeyId:
			errstr = awsErr.Error()
		case ssm.ErrCodeParameterNotFound:
			errstr = awsErr.Error()
		case ssm.ErrCodeParameterVersionNotFound:
			errstr = awsErr.Error()
		}
	}
	fmt.Printf("Error: %s", errstr)
}

func CheckSSMByPath(err error) {
	var errstr string
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ssm.ErrCodeInternalServerError:
			errstr = awsErr.Error()
		case ssm.ErrCodeInvalidFilterKey:
			errstr = awsErr.Error()
		case ssm.ErrCodeInvalidFilterOption:
			errstr = awsErr.Error()
		case ssm.ErrCodeInvalidFilterValue:
			errstr = awsErr.Error()
		case ssm.ErrCodeInvalidKeyId:
			errstr = awsErr.Error()
		case ssm.ErrCodeInvalidNextToken:
			errstr = awsErr.Error()
		}
	}
	fmt.Printf("SSM By Path Error: %s", errstr)
}
