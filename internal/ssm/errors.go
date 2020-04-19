package ssm

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
)

// checkSSMError is the entry point for check all of the based and call specific errors
func checkSSMError(a action, err error) error {
	awsErr := checkBaseSSMErrors(err)
	switch a {
	case Get:
		return checkSSMGetParameterError(err)
	case Save:
		return checkSSMPutParameterError(err)
	}
	return awsErr
}

// checkRegion catches if the aws region isn't configured
func checkRegion(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch  awsErr.Code() {
		case "MissingRegion":
			return awsErr
		}
	}
	return nil
}

// checkBaseSSMErrors checks for the common errors all SSM API calls might return
func checkBaseSSMErrors(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInvalidKeyId:
			return awsErr
		case awsssm.ErrCodeInternalServerError:
			return awsErr
		}
	}
	return checkRegion(err)
}

// checkSSMGetParameterError checks for errors the GetParameter API call might return
func checkSSMGetParameterError(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeParameterNotFound:
			return awsErr
		default:
			// this means we are missing a check and we can print it out in the calling function
			return awsErr
		}
	}
	//return errors.New("SSM Get Param error")
	return nil
}

// checkSSMPutParameterError checks for errors the PutParameter API call might return
func checkSSMPutParameterError(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		//case awsssm.ErrCodeInvalidPa:
		//	return awsErr
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
		case awsssm.ErrCodeParameterAlreadyExists:
			return awsErr
		default:
			// this means we are missing a check and we can print it out in the calling function
			return awsErr
		}
	}
	//return errors.New("aws ssm param put error")
	return nil
}

// checkSSMByPathError checks for errors the GetParameterByPath API call might return
func checkSSMByPathError(err error) error {
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
	//return errors.New("aws ssm param get by path error")
	return nil
}
