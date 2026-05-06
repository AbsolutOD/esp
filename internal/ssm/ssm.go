package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pinpt/esp/internal/common"
	"github.com/pinpt/esp/internal/utils"
)

type action string

// Constants to represent actions to take against SSM
const (
	Get     action = "get"
	GetMany action = "getMany"
	//Put     action = "put"
	Save   action = "save"
	Delete action = "delete"
)

// Service the struct representing AWS service/Session
type Service struct {
	Svc    *awsssm.Client
	Region string
	Cfg    aws.Config
}

// mapErr applies the per-action error mapper. If the mapper returns
// nil (the error wasn't a recognized AWS error type), return the raw
// error so the caller still sees the failure.
func mapErr(a action, err error) error {
	if mapped := checkSSMError(a, err); mapped != nil {
		return mapped
	}
	return err
}

// Save a single param for a given path
func (s *Service) Save(p common.EspParamInput) (common.SaveOutput, error) {
	pi := &awsssm.PutParameterInput{
		Type:  selectType(p.Secure),
		Name:  aws.String(p.Name),
		Value: aws.String(p.Value),
	}
	param, err := s.Svc.PutParameter(context.Background(), pi)
	if err != nil {
		return common.SaveOutput{}, mapErr(Save, err)
	}
	return common.SaveOutput{Version: param.Version}, nil
}

// Delete a single param for a given path
func (s *Service) Delete(p common.DeleteInput) (string, error) {
	dpi := &awsssm.DeleteParameterInput{
		Name: aws.String(p.Name),
	}
	_, err := s.Svc.DeleteParameter(context.Background(), dpi)
	if err != nil {
		return "", mapErr(Delete, err)
	}
	return p.Name, nil
}

// GetOne gets a single param for a given path
func (s *Service) GetOne(p common.GetOneInput) (common.EspParam, error) {
	si := &awsssm.GetParameterInput{
		Name:           aws.String(p.Name),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	resp, err := s.Svc.GetParameter(context.Background(), si)
	if err != nil {
		return common.EspParam{}, mapErr(Get, err)
	}
	return convertToEspParam(*resp.Parameter), nil
}

// GetMany recursively gets parameters from a given path
func (s *Service) GetMany(p common.ListParamInput) ([]common.EspParam, error) {
	si := &awsssm.GetParametersByPathInput{
		Path:           aws.String(p.Path),
		WithDecryption: aws.Bool(p.Decrypt),
		Recursive:      aws.Bool(p.Recursive),
	}
	paginator := awsssm.NewGetParametersByPathPaginator(s.Svc, si)

	var espParams []common.EspParam
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, mapErr(GetMany, err)
		}
		for _, v := range page.Parameters {
			espParams = append(espParams, convertToEspParam(v))
		}
	}
	return espParams, nil
}

// Copy method copies the given parameter to a new location
func (s *Service) Copy(cc common.CopyCommand) (common.SaveOutput, error) {
	input := common.GetOneInput{
		Name:    cc.Source,
		Decrypt: true,
	}
	sparam, err := s.GetOne(input)
	if err != nil {
		return common.SaveOutput{}, err
	}

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
	return svc
}

// Init create the actual session to talk to the AWS API
func (s *Service) Init() error {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(s.Region))
	if err != nil {
		return err
	}
	s.Cfg = cfg
	s.Svc = awsssm.NewFromConfig(cfg)
	return nil
}
