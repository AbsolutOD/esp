package ssm

import (
	"errors"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pinpt/esp/internal/client"
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

func ConvertToEspParam(ap *awsssm.GetParameterOutput) client.EspParam {

}