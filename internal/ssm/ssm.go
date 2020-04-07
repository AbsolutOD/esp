package ssm

import (
	"fmt"
	"github.com/pinpt/esp/internal/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pinpt/esp/internal/common"
	"os"
)

type action string

const (
	Get action = "get"
	GetMany action = "getMany"
	Put action = "put"
	Save action = "save"
)

type Service struct {
	Svc *awsssm.SSM
	Region string
	Cfg aws.Config
	session *session.Session
}

func (s *Service) Save(p common.EspParamInput) common.SaveOutput {
	pi := &awsssm.PutParameterInput{
		Type: SelectType(p.Secure),
		Name: aws.String(p.Name),
		Value: aws.String(p.Value),
	}
	param, err := s.Svc.PutParameter(pi)
	if err != nil {
		awsErr := CheckSSMError(Save, err)
		fmt.Println("SSM Error: %s", awsErr.Error())
		os.Exit(1)
	}
	return common.SaveOutput{ Version: *param.Version }
}

func (s *Service) GetOne(p common.GetOneInput) common.EspParam {
	si := &awsssm.GetParameterInput{
		Name: aws.String(p.Name),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	resp, err := s.Svc.GetParameter(si)
	if err != nil {
		awsErr := CheckSSMError(Get, err)
		fmt.Println("SSM Error: %s", awsErr.Error())
		os.Exit(1)
	}
	param := ConvertToEspParam(resp.Parameter)
	return param
}

func (s *Service) GetMany(p common.ListParamInput) []common.EspParam {
	si := &awsssm.GetParametersByPathInput{
		Path:           aws.String(p.Path),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	params, err := s.Svc.GetParametersByPath(si)
	if err != nil {
		awsErr := CheckSSMError(GetMany, err)
		fmt.Println("SSM Error: %s", awsErr.Error())
		os.Exit(1)
	}

	var espParams []common.EspParam
	for _, v := range params.Parameters {
		espParams = append(espParams, ConvertToEspParam(v))
	}

	return espParams
}

// New Create a new SSM service
func New() *Service {
	svc := new(Service)
	svc.Region = utils.GetEnv("AWS_REGION", "us-east-1")
	svc.Cfg = aws.Config{ Region: aws.String(svc.Region) }
	return svc
}

// Init create the actual session to talk to the AWS API
func (s *Service) Init()  {
	s.session = session.Must(session.NewSession(&s.Cfg))
	s.Svc = awsssm.New(s.session)
}
