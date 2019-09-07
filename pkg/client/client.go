package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type EspClient struct {
	Client ssmiface.SSMAPI
	Region string
	Cfg aws.Config
	//client ssm.Session
}

func (e EspClient) New(region string) {
	e.Region = region
	e.Cfg = aws.Config{
		Region: aws.String(region),
	}
}

func (e EspClient) GetSsmClient(region string) *ssm.SSM {
	sess := session.Must(session.NewSession(&e.Cfg))
	e.Client = ssm.New(sess)
}


