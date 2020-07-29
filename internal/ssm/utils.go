package ssm

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pinpt/esp/internal/common"
)

// ParamType sets the base type for SSM parameter types
type ParamType string

// Defines the SSM types
const (
	String       ParamType = "string"
	SecureString ParamType = "SecureString"
	StringList   ParamType = "Stringlist"
)

// AwsParam represents an individual SSM parameter
type AwsParam struct {
	Arn              string
	Name             string
	Type             ParamType
	Value            string
	Version          int
	LastModifiedDate float32
}

func (p *AwsParam) isValid() error {
	switch p.Type {
	case String, SecureString, StringList:
		return nil
	}
	return errors.New("invalid SSM Parameter Type")
}

func selectType(t bool) *string {
	if t {
		return aws.String(awsssm.ParameterTypeSecureString)
	}

	return aws.String(awsssm.ParameterTypeString)
}

func convertToEspParam(ap *awsssm.Parameter) common.EspParam {
	param := common.EspParam{
		Id:               *ap.ARN,
		Name:             *ap.Name,
		Type:             *ap.Type,
		Value:            *ap.Value,
		Version:          *ap.Version,
		LastModifiedDate: *ap.LastModifiedDate,
	}

	if param.Type == "SecureString" {
		param.Secure = true
	}
	return param
}

// handleAwsErr it will for all of the AWS API errors and exit if exists
func handleAwsErr(a action, err error) {
	awsErr := checkSSMError(a, err)
	if awsErr != nil {
		fmt.Printf("SSM Error: %s\n", awsErr.Error())
		os.Exit(1)
	}
}
