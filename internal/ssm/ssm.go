package ssm

import (
	"github.com/absolutod/esp/internal/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/absolutod/esp/internal/common"
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
		CheckSSMError(Save, err)
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
		CheckSSMError(Get, err)
	}
	param := ConvertToEspParam(resp)
	return param
}

/*func GetMany(ec common.EspConfig, d bool, paths []*string) []*ssm.Parameter {

	si := &ssm.GetParametersInput{
		Names: paths,
		WithDecryption: aws.Bool(d),
	}
	fmt.Println(si)
	resp, err := ec.Svc.GetParameters(si)
	if err != nil {
		fmt.Println(err)
		errors.CheckSSMGetParameters(err)
	}
	return resp.Parameters
}*/

// actually create the ssm common
func New() *Service {
	svc := new(Service)
	svc.Region = utils.GetEnv("AWS_REGION", "us-east-1")
	svc.Cfg = aws.Config{ Region: aws.String(svc.Region) }
	svc.session = session.Must(session.NewSession(&svc.Cfg))
	svc.Svc = awsssm.New(svc.session)
	return svc
}
