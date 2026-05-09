package ssm

import (
	"context"

	"github.com/AbsolutOD/esp/internal/common"
	"github.com/AbsolutOD/esp/internal/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
)

type action string

const (
	Get     action = "get"
	GetMany action = "getMany"
	Save    action = "save"
	Delete  action = "delete"
)

// ssmAPI is the subset of the AWS SSM client used by Service.
// The concrete *awsssm.Client satisfies it; tests inject a fake.
// The four-method shape is required by NewGetParametersByPathPaginator,
// whose first argument is awsssm.GetParametersByPathAPIClient.
type ssmAPI interface {
	PutParameter(ctx context.Context, in *awsssm.PutParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.PutParameterOutput, error)
	GetParameter(ctx context.Context, in *awsssm.GetParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.GetParameterOutput, error)
	DeleteParameter(ctx context.Context, in *awsssm.DeleteParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.DeleteParameterOutput, error)
	GetParametersByPath(ctx context.Context, in *awsssm.GetParametersByPathInput, optFns ...func(*awsssm.Options)) (*awsssm.GetParametersByPathOutput, error)
}

// Service is the SSM-backed implementation of client.Client.
type Service struct {
	api    ssmAPI
	Region string
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

// New builds a Service backed by a real AWS SSM client. Returns an
// error if AWS config loading fails.
func New() (*Service, error) {
	region := utils.GetEnv("AWS_REGION", "us-east-1")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &Service{api: awsssm.NewFromConfig(cfg), Region: region}, nil
}

// Save a single param for a given path
func (s *Service) Save(p common.EspParamInput) (common.SaveOutput, error) {
	pi := &awsssm.PutParameterInput{
		Type:  selectType(p.Secure),
		Name:  &p.Name,
		Value: &p.Value,
	}
	param, err := s.api.PutParameter(context.Background(), pi)
	if err != nil {
		return common.SaveOutput{}, mapErr(Save, err)
	}
	return common.SaveOutput{Version: param.Version}, nil
}

// Delete a single param for a given path
func (s *Service) Delete(p common.DeleteInput) (string, error) {
	dpi := &awsssm.DeleteParameterInput{
		Name: &p.Name,
	}
	_, err := s.api.DeleteParameter(context.Background(), dpi)
	if err != nil {
		return "", mapErr(Delete, err)
	}
	return p.Name, nil
}

// GetOne gets a single param for a given path
func (s *Service) GetOne(p common.GetOneInput) (common.EspParam, error) {
	si := &awsssm.GetParameterInput{
		Name:           &p.Name,
		WithDecryption: &p.Decrypt,
	}
	resp, err := s.api.GetParameter(context.Background(), si)
	if err != nil {
		return common.EspParam{}, mapErr(Get, err)
	}
	return convertToEspParam(*resp.Parameter), nil
}

// GetMany recursively gets parameters from a given path
func (s *Service) GetMany(p common.ListParamInput) ([]common.EspParam, error) {
	si := &awsssm.GetParametersByPathInput{
		Path:           &p.Path,
		WithDecryption: &p.Decrypt,
		Recursive:      &p.Recursive,
	}
	paginator := awsssm.NewGetParametersByPathPaginator(s.api, si)

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
	sparam, err := s.GetOne(common.GetOneInput{Name: cc.Source, Decrypt: true})
	if err != nil {
		return common.SaveOutput{}, err
	}
	return s.Save(common.EspParamInput{
		Name:   cc.Destination,
		Secure: sparam.Secure,
		Value:  sparam.Value,
	})
}
