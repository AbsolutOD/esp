package ssmparam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pinpt/esp/pkg/client"
	"github.com/pinpt/esp/pkg/errors"
)

// getParam Queries the ssm param
func GetOne(ec client.EspConfig, d bool, path string) *ssm.Parameter {
	si := &ssm.GetParameterInput{
		Name: aws.String(path),
		WithDecryption: aws.Bool(d),
	}
	resp, err := ec.Svc.GetParameter(si)
	if err != nil {
		errors.CheckSSMGetParameters(err)
	}

	return resp.Parameter
}

/*func GetMany(ec client.EspConfig, d bool, paths []*string) []*ssm.Parameter {

	si := &ssm.GetParametersInput{
		Names: &paths,
		WithDecryption: aws.Bool(d),
	}
	return
}*/
