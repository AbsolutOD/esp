package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pinpt/esp/internal/client"
	"github.com/absolutod/esp/internal/errors"
	"os"
)

type Service struct {
	Svc *awsssm.SSM
	Region string
	Cfg aws.Config
	session *session.Session
}

func (s *Service) Save(p client.EspParam) (client.EspParam, error) {
	panic("implement me")
}

func (s *Service) GetOne(p client.GetOneInput) (client.EspParam, error) {
	si := &awsssm.GetParameterInput{
		Name: aws.String(p.Path),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	resp, err := s.Svc.GetParameter(si)
	if err != nil {
		CheckSSMGetParameters(err)
		return client.EspParam{}, errors.New("")
	}
	param := ConvertToEspParam(resp)
	return param, nil
}

/*func GetMany(ec client.EspConfig, d bool, paths []*string) []*ssm.Parameter {

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

// actually create the ssm client
func New() *Service {
	svc := new(Service)
	svc.Region = os.Getenv("AWS_REGION")
	svc.Cfg = aws.Config{ Region: aws.String(svc.Region) }
	svc.session = session.Must(session.NewSession(&s.Cfg))
	svc.Svc = awsssm.New(s.session)
	return svc
}

