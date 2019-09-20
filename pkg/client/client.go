package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type EspConfig struct {
	Client ssmiface.SSMAPI
	Region string
	Cfg aws.Config
	//client ssm.Session
}

func New(region string) EspConfig {
	e := EspConfig{
		Region: region,
		Cfg: aws.Config{
			Region: aws.String(region),
		},
	}
	return e
}

func (e EspConfig) GetSsmClient(region string) {
	sess := session.Must(session.NewSession(&e.Cfg))
	e.Client = ssm.New(sess)
}
