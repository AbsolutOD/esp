package ssm

import (
	"errors"
	"github.com/pinpt/esp/internal/common"
	"github.com/aws/aws-sdk-go/aws"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
)

type ParamType string

const(
	String ParamType = "string"
	SecureString ParamType = "securestring"
	StringList ParamType = "stringlist"
)

type AwsParam struct {
	Arn string
	Name string
	Type ParamType
	Value string
	Version int
	LastModifiedDate float32
}

func (p *AwsParam) IsValid() error {
	switch p.Type {
	case String, SecureString, StringList:
		return nil
	}
	return errors.New("invalid SSM Parameter Type")
}

func SelectType(t bool) *string {
	if t {
		return aws.String(awsssm.ParameterTypeSecureString)
	} else {
		return aws.String(awsssm.ParameterTypeString)
	}
}

func ConvertToEspParam(ap *awsssm.Parameter) common.EspParam {
	param := common.EspParam{
		Id: *ap.ARN,
		Name: *ap.Name,
		Type: *ap.Type,
		Value: *ap.Value,
		Version: *ap.Version,
		LastModifiedDate: *ap.LastModifiedDate,
	}

	if param.Type == "securestring" {
		param.Secure = true
	}
	return param
}