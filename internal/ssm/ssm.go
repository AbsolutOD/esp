package ssm

import (
	"context"
	"fmt"
	"os"

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

// Save a single param for a given path
func (s *Service) Save(p common.EspParamInput) common.SaveOutput {
	pi := &awsssm.PutParameterInput{
		Type:  selectType(p.Secure),
		Name:  aws.String(p.Name),
		Value: aws.String(p.Value),
	}
	param, err := s.Svc.PutParameter(context.Background(), pi)
	if err != nil {
		handleAwsErr(Save, err)
	}
	return common.SaveOutput{Version: param.Version}
}

// Delete a single param for a given path
func (s *Service) Delete(p common.DeleteInput) string {
	dpi := &awsssm.DeleteParameterInput{
		Name: aws.String(p.Name),
	}
	_, err := s.Svc.DeleteParameter(context.Background(), dpi)
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
	resp, err := s.Svc.GetParameter(context.Background(), si)
	if err != nil {
		handleAwsErr(Get, err)
	}
	return convertToEspParam(*resp.Parameter)
}

// GetMany recursively gets parameters from a given path
func (s *Service) GetMany(p common.ListParamInput) []common.EspParam {
	si := &awsssm.GetParametersByPathInput{
		Path:           aws.String(p.Path),
		WithDecryption: aws.Bool(p.Decrypt),
		Recursive:      aws.Bool(p.Recursive),
	}
	params, err := s.Svc.GetParametersByPath(context.Background(), si)
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
	params, err := s.Svc.GetParametersByPath(context.Background(), pi)
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
	return svc
}

// Init create the actual session to talk to the AWS API
func (s *Service) Init() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(s.Region))
	if err != nil {
		fmt.Printf("AWS config load error: %s\n", err.Error())
		os.Exit(1)
	}
	s.Cfg = cfg
	s.Svc = awsssm.NewFromConfig(cfg)
}
