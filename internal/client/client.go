package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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
