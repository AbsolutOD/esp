package ssm

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/AbsolutOD/esp/internal/common"
)

func selectType(t bool) ssmtypes.ParameterType {
	if t {
		return ssmtypes.ParameterTypeSecureString
	}
	return ssmtypes.ParameterTypeString
}

func convertToEspParam(ap ssmtypes.Parameter) common.EspParam {
	param := common.EspParam{
		Id:               aws.ToString(ap.ARN),
		Name:             aws.ToString(ap.Name),
		Type:             string(ap.Type),
		Value:            aws.ToString(ap.Value),
		Version:          ap.Version,
		LastModifiedDate: aws.ToTime(ap.LastModifiedDate),
	}

	if param.Type == string(ssmtypes.ParameterTypeSecureString) {
		param.Secure = true
	}
	return param
}
