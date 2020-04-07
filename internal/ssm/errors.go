package ssm

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
)

// CheckSSMError is the entry point for check all of the based and call specific errors
func CheckSSMError(a action, err error) error {
	awsErr := CheckBaseSSMErrors(err)
	if awsErr != nil {
		return awsErr
	}
	switch a {
	case Get:
		return CheckSSMGetParameterError(err)
	case Save:
		return CheckSSMPutParameterError(err)
	}
}

// CheckRegion catches if the aws region isn't configured
func CheckRegion(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch  awsErr.Code() {
		case "MissingRegion":
			return awsErr
		}
	}
}

// CheckBaseSSMErrors checks for the common errors all SSM API calls might return
func CheckBaseSSMErrors(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInvalidKeyId:
			return awsErr
		case awsssm.ErrCodeInternalServerError:
			return awsErr
		default:
			return CheckRegion(err)
		}
	}
}

// CheckSSMGetParameterError checks for errors the GetParameter API call might return
func CheckSSMGetParameterError(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeParameterNotFound:
			return awsErr
		}
	}
}

// CheckSSMPutParameterError checks for errors the PutParameter API call might return
func CheckSSMPutParameterError(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeParameterLimitExceeded:
			return awsErr
		case awsssm.ErrCodeTooManyUpdates:
			return awsErr
		case awsssm.ErrCodeHierarchyTypeMismatchException:
			return awsErr
		case awsssm.ErrCodeInvalidAllowedPatternException:
			return awsErr
		case awsssm.ErrCodeParameterMaxVersionLimitExceeded:
			return awsErr
		case awsssm.ErrCodeUnsupportedParameterType:
			return awsErr
		case awsssm.ErrCodePoliciesLimitExceededException:
			return awsErr
		case awsssm.ErrCodeInvalidPolicyTypeException:
			return awsErr
		case awsssm.ErrCodeInvalidPolicyAttributeException:
			return awsErr
		case awsssm.ErrCodeIncompatiblePolicyException:
			return awsErr
		}
	}
}

// CheckSSMByPathError checks for errors the GetParameterByPath API call might return
func CheckSSMByPathError(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInternalServerError:
			return awsErr
		case awsssm.ErrCodeInvalidFilterKey:
			return awsErr
		case awsssm.ErrCodeInvalidFilterOption:
			return awsErr
		case awsssm.ErrCodeInvalidFilterValue:
			return awsErr
		case awsssm.ErrCodeInvalidKeyId:
			return awsErr
		case awsssm.ErrCodeInvalidNextToken:
			return awsErr
		}
	}
}
