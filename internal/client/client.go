package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pinpt/esp/internal/errors"
)

type EspConfig struct {
	Svc *ssm.SSM
	//Client ssmiface.SSMAPI
	Region string
	Cfg aws.Config
	session *session.Session
}

func New(region string) EspConfig {
	e := EspConfig{
		Region: region,
		Cfg: aws.Config{
			Region: aws.String(region),
		},
	}
	e.GetSsmClient()
	return e
}

// actually create the ssm client
func (e *EspConfig) GetSsmClient() {
	e.session = session.Must(session.NewSession(&e.Cfg))
	e.Svc = ssm.New(e.session)
}

// getParam Queries the ssm param
func (e *EspConfig) getParam(d bool, key string) *ssm.Parameter {
	si := &ssm.GetParameterInput{
		Name: aws.String(key),
		WithDecryption: aws.Bool(d),
	}
	resp, err := e.Svc.GetParameter(si)
	if err != nil {
		errors.CheckSSMGetParameters(err)
	}

	return resp.Parameter
}

