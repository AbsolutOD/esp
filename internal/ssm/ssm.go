package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pinpt/esp/internal/common"
	"github.com/pinpt/esp/internal/utils"
)

type action string

// Constants to represent actions to take against SSM
const (
	Get     action = "get"
	GetMany action = "getMany"
	Put     action = "put"
	Save    action = "save"
	Delete  action = "delete"
)

// Service the struct representing AWS service/Session
type Service struct {
	Svc     *awsssm.SSM
	Region  string
	Cfg     aws.Config
	session *session.Session
}

// Save a single param for a given path
func (s *Service) Save(p common.EspParamInput) common.SaveOutput {
	pi := &awsssm.PutParameterInput{
		Type:  selectType(p.Secure),
		Name:  aws.String(p.Name),
		Value: aws.String(p.Value),
	}
	param, err := s.Svc.PutParameter(pi)
	if err != nil {
		handleAwsErr(Save, err)
	}
	return common.SaveOutput{Version: *param.Version}
}

// Delete a single param for a given path
func (s *Service) Delete(p common.DeleteInput) string {
	dpi := &awsssm.DeleteParameterInput{
		Name: aws.String(p.Name),
	}
	_, err := s.Svc.DeleteParameter(dpi)
	if err != nil {
		handleAwsErr(Delete, err)
	}
	return p.Name
}

// GetOne gets a single param for a given path
func (s *Service) GetOne(p common.GetOneInput) common.EspParam {
	si := &awsssm.GetParameterInput{
		Name:           aws.String(p.Name),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	resp, err := s.Svc.GetParameter(si)
	if err != nil {
		handleAwsErr(Get, err)
	}
	param := convertToEspParam(resp.Parameter)
	return param
}

// GetMany recursively gets parameters from a given path
func (s *Service) GetMany(p common.ListParamInput) []common.EspParam {
	si := &awsssm.GetParametersByPathInput{
		Path:           aws.String(p.Path),
		WithDecryption: aws.Bool(p.Decrypt),
		Recursive:      aws.Bool(p.Recursive),
	}
	params, err := s.Svc.GetParametersByPath(si)
	if err != nil {
		handleAwsErr(GetMany, err)
	}

	var espParams []common.EspParam
	for _, v := range params.Parameters {
		espParams = append(espParams, convertToEspParam(v))
	}

	if params.NextToken != nil {
		si.NextToken = params.NextToken
		moreParams := s.getNextParams(si)
		espParams = append(espParams, moreParams...)
	}
	return espParams
}

// getNextParams uses the NextToken to get more params
func (s *Service) getNextParams(pi *awsssm.GetParametersByPathInput) []common.EspParam {
	params, err := s.Svc.GetParametersByPath(pi)
	if err != nil {
		handleAwsErr(GetMany, err)
	}

	var espParams []common.EspParam
	for _, v := range params.Parameters {
		espParams = append(espParams, convertToEspParam(v))
	}

	if params.NextToken != nil {
		pi.NextToken = params.NextToken
		moreParams := s.getNextParams(pi)
		espParams = append(espParams, moreParams...)
	}
	return espParams
}

// Copy method copies the given parameter to a new location
func (s *Service) Copy(cc common.CopyCommand) common.SaveOutput {
	input := common.GetOneInput{
		Name:    cc.Source,
		Decrypt: true,
	}
	sparam := s.GetOne(input)

	dparam := common.EspParamInput{
		Name:   cc.Destination,
		Secure: sparam.Secure,
		Value:  sparam.Value,
	}
	return s.Save(dparam)
}

// New Create a new SSM service
func New() *Service {
	svc := new(Service)
	svc.Region = utils.GetEnv("AWS_REGION", "us-east-1")
	svc.Cfg = aws.Config{Region: aws.String(svc.Region)}
	return svc
}

// Init create the actual session to talk to the AWS API
func (s *Service) Init() {
	s.session = session.Must(session.NewSession(&s.Cfg))
	s.Svc = awsssm.New(s.session)
}
