package ssm

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"os"
)

func CheckSSMError(a action, err error)  {
	CheckBaseSSMErrors(err)
	switch a {
	case Get:
		CheckSSMGetParameterError(err)
	case Save:
		CheckSSMPutParameterError(err)
	}
}

func CheckRegion(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		switch  awsErr.Code() {
		case "MissingRegion":
			fmt.Println(aws.ErrMissingRegion)
			os.Exit(1)
		}
	}
}

func CheckBaseSSMErrors(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInvalidKeyId:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInternalServerError:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		default:
			CheckRegion(err)
		}
	}
}
func CheckSSMGetParameterError(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeParameterNotFound:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		default:
			CheckRegion(err)
		}
	}
}

func CheckSSMPutParameterError(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeParameterLimitExceeded:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeTooManyUpdates:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeHierarchyTypeMismatchException:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidAllowedPatternException:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeParameterMaxVersionLimitExceeded:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeUnsupportedParameterType:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodePoliciesLimitExceededException:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidPolicyTypeException:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidPolicyAttributeException:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeIncompatiblePolicyException:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		}
	}
}

func CheckSSMByPathError(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case awsssm.ErrCodeInternalServerError:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidFilterKey:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidFilterOption:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidFilterValue:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidKeyId:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		case awsssm.ErrCodeInvalidNextToken:
			fmt.Println("SSM Error: %s", awsErr.Error())
			os.Exit(1)
		}
	}
}
