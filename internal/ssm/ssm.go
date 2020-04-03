package ssm

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/absolutod/esp/internal/common"
	"os"
)

type Service struct {
	Svc *awsssm.SSM
	Region string
	Cfg aws.Config
	session *session.Session
}

func (s *Service) Save(p common.EspParam) (common.EspParam, error) {
	panic("implement me")
}

func (s *Service) GetOne(p common.GetOneInput) (common.EspParam, error) {
	si := &awsssm.GetParameterInput{
		Name: aws.String(p.Name),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	resp, err := s.Svc.GetParameter(si)
	if err != nil {
		//return common.EspParam{}, errors.New("fubar error")
		CheckSSMGetParameters(err)
		return common.EspParam{}, errors.New("Error Getting the ssm parameter")
	}
	param := ConvertToEspParam(resp)
	return param, nil
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
	svc.Region = os.Getenv("AWS_REGION")
	svc.Cfg = aws.Config{ Region: aws.String(svc.Region) }
	svc.session = session.Must(session.NewSession(&svc.Cfg))
	svc.Svc = awsssm.New(svc.session)
	return svc
}
