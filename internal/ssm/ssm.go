package ssm

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/absolutod/esp/internal/client"
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
		Name: aws.String(p.Name),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	resp, err := s.Svc.GetParameter(si)
	if err != nil {
		return client.EspParam{}, errors.New("fubar error")
		//CheckSSMGetParameters(err)
		//return client.EspParam{}, errors.New("Error Getting the ssm parameter")
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
	svc.session = session.Must(session.NewSession(&svc.Cfg))
	svc.Svc = awsssm.New(svc.session)
	return svc
}
