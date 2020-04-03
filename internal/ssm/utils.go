package ssm

import (
	"errors"
	"github.com/absolutod/esp/internal/common"
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

func ConvertToEspParam(ap *awsssm.GetParameterOutput) common.EspParam {
	param := common.EspParam{
		Id: *ap.Parameter.ARN,
		Name: *ap.Parameter.Name,
		Type: *ap.Parameter.Type,
		Value: *ap.Parameter.Value,
		Version: *ap.Parameter.Version,
		LastModifiedDate: *ap.Parameter.LastModifiedDate,
	}

	if param.Type == "securestring" {
		param.Secure = true
	}
	return param
}