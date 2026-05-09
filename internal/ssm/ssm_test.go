package ssm

import (
	"context"
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/common"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// fakeSSMAPI implements ssmAPI by returning canned responses.
// Each method records its last input for assertion.
type fakeSSMAPI struct {
	putIn    *awsssm.PutParameterInput
	putOut   *awsssm.PutParameterOutput
	putErr   error
	getIn    *awsssm.GetParameterInput
	getOut   *awsssm.GetParameterOutput
	getErr   error
	delIn    *awsssm.DeleteParameterInput
	delOut   *awsssm.DeleteParameterOutput
	delErr   error
	pathIn   *awsssm.GetParametersByPathInput
	pathOuts []*awsssm.GetParametersByPathOutput
	pathErrs []error
	pathIdx  int
}

func (f *fakeSSMAPI) PutParameter(_ context.Context, in *awsssm.PutParameterInput, _ ...func(*awsssm.Options)) (*awsssm.PutParameterOutput, error) {
	f.putIn = in
	return f.putOut, f.putErr
}

func (f *fakeSSMAPI) GetParameter(_ context.Context, in *awsssm.GetParameterInput, _ ...func(*awsssm.Options)) (*awsssm.GetParameterOutput, error) {
	f.getIn = in
	return f.getOut, f.getErr
}

func (f *fakeSSMAPI) DeleteParameter(_ context.Context, in *awsssm.DeleteParameterInput, _ ...func(*awsssm.Options)) (*awsssm.DeleteParameterOutput, error) {
	f.delIn = in
	return f.delOut, f.delErr
}

func (f *fakeSSMAPI) GetParametersByPath(_ context.Context, in *awsssm.GetParametersByPathInput, _ ...func(*awsssm.Options)) (*awsssm.GetParametersByPathOutput, error) {
	f.pathIn = in
	i := f.pathIdx
	f.pathIdx++
	if i >= len(f.pathOuts) {
		return &awsssm.GetParametersByPathOutput{}, nil
	}
	var err error
	if i < len(f.pathErrs) {
		err = f.pathErrs[i]
	}
	return f.pathOuts[i], err
}

func newServiceWithAPI(api ssmAPI) *Service {
	return &Service{api: api, Region: "us-east-1"}
}

func strPtr(s string) *string { return &s }

// --- Save ---

func TestService_Save_Success(t *testing.T) {
	fake := &fakeSSMAPI{putOut: &awsssm.PutParameterOutput{Version: 7}}
	s := newServiceWithAPI(fake)

	out, err := s.Save(common.EspParamInput{Name: "/x/y", Value: "v", Secure: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Version != 7 {
		t.Errorf("Version = %d, want 7", out.Version)
	}
	if fake.putIn.Type != ssmtypes.ParameterTypeString {
		t.Errorf("Type = %q, want String", fake.putIn.Type)
	}
	if got := *fake.putIn.Name; got != "/x/y" {
		t.Errorf("Name = %q, want /x/y", got)
	}
	if got := *fake.putIn.Value; got != "v" {
		t.Errorf("Value = %q, want v", got)
	}
}

func TestService_Save_SecureFlagSetsSecureString(t *testing.T) {
	fake := &fakeSSMAPI{putOut: &awsssm.PutParameterOutput{Version: 1}}
	s := newServiceWithAPI(fake)
	if _, err := s.Save(common.EspParamInput{Name: "n", Value: "v", Secure: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.putIn.Type != ssmtypes.ParameterTypeSecureString {
		t.Errorf("Type = %q, want SecureString", fake.putIn.Type)
	}
}

func TestService_Save_AlreadyExistsErrorPropagates(t *testing.T) {
	awsErr := &ssmtypes.ParameterAlreadyExists{}
	fake := &fakeSSMAPI{putErr: awsErr}
	s := newServiceWithAPI(fake)

	_, err := s.Save(common.EspParamInput{Name: "n", Value: "v"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterAlreadyExists
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterAlreadyExists", err)
	}
}

// --- GetOne ---

func TestService_GetOne_Success(t *testing.T) {
	fake := &fakeSSMAPI{getOut: &awsssm.GetParameterOutput{
		Parameter: &ssmtypes.Parameter{
			Name:    strPtr("/x/y"),
			Value:   strPtr("v"),
			Type:    ssmtypes.ParameterTypeString,
			Version: 1,
		},
	}}
	s := newServiceWithAPI(fake)

	got, err := s.GetOne(common.GetOneInput{Name: "/x/y", Decrypt: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "/x/y" {
		t.Errorf("Name = %q, want /x/y", got.Name)
	}
	if got.Value != "v" {
		t.Errorf("Value = %q, want v", got.Value)
	}
	if !*fake.getIn.WithDecryption {
		t.Error("WithDecryption = false, want true")
	}
}

func TestService_GetOne_NotFoundErrorPropagates(t *testing.T) {
	fake := &fakeSSMAPI{getErr: &ssmtypes.ParameterNotFound{}}
	s := newServiceWithAPI(fake)

	_, err := s.GetOne(common.GetOneInput{Name: "/x/y"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterNotFound
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterNotFound", err)
	}
}

// --- GetMany ---

func TestService_GetMany_MultiplePages(t *testing.T) {
	fake := &fakeSSMAPI{
		pathOuts: []*awsssm.GetParametersByPathOutput{
			{
				Parameters: []ssmtypes.Parameter{
					{Name: strPtr("/a/1"), Value: strPtr("v1"), Type: ssmtypes.ParameterTypeString},
				},
				NextToken: strPtr("page2"),
			},
			{
				Parameters: []ssmtypes.Parameter{
					{Name: strPtr("/a/2"), Value: strPtr("v2"), Type: ssmtypes.ParameterTypeString},
				},
			},
		},
	}
	s := newServiceWithAPI(fake)

	params, err := s.GetMany(common.ListParamInput{Path: "/a/", Recursive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(params) != 2 {
		t.Fatalf("len(params) = %d, want 2", len(params))
	}
	if params[0].Name != "/a/1" || params[1].Name != "/a/2" {
		t.Errorf("got names %q,%q; want /a/1,/a/2", params[0].Name, params[1].Name)
	}
}

func TestService_GetMany_ErrorMidIteration(t *testing.T) {
	fake := &fakeSSMAPI{
		pathOuts: []*awsssm.GetParametersByPathOutput{
			{
				Parameters: []ssmtypes.Parameter{{Name: strPtr("/a/1"), Type: ssmtypes.ParameterTypeString}},
				NextToken:  strPtr("page2"),
			},
			{},
		},
		pathErrs: []error{nil, &ssmtypes.InternalServerError{}},
	}
	s := newServiceWithAPI(fake)

	_, err := s.GetMany(common.ListParamInput{Path: "/a/"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.InternalServerError
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *InternalServerError", err)
	}
}

func TestService_GetMany_EmptyPages(t *testing.T) {
	fake := &fakeSSMAPI{
		pathOuts: []*awsssm.GetParametersByPathOutput{{}},
	}
	s := newServiceWithAPI(fake)
	params, err := s.GetMany(common.ListParamInput{Path: "/x/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(params) != 0 {
		t.Errorf("got %d params, want 0", len(params))
	}
}

// --- Delete ---

func TestService_Delete_Success(t *testing.T) {
	fake := &fakeSSMAPI{delOut: &awsssm.DeleteParameterOutput{}}
	s := newServiceWithAPI(fake)

	name, err := s.Delete(common.DeleteInput{Name: "/x/y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "/x/y" {
		t.Errorf("name = %q, want /x/y", name)
	}
	if got := *fake.delIn.Name; got != "/x/y" {
		t.Errorf("input name = %q, want /x/y", got)
	}
}

func TestService_Delete_NotFoundErrorPropagates(t *testing.T) {
	fake := &fakeSSMAPI{delErr: &ssmtypes.ParameterNotFound{}}
	s := newServiceWithAPI(fake)

	_, err := s.Delete(common.DeleteInput{Name: "/x/y"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterNotFound
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterNotFound", err)
	}
}

// --- Copy ---

func TestService_Copy_Success(t *testing.T) {
	// First call: GetOne returns the source param.
	// Second call: PutParameter.
	fake := &fakeSSMAPI{
		getOut: &awsssm.GetParameterOutput{
			Parameter: &ssmtypes.Parameter{
				Name:  strPtr("/src"),
				Value: strPtr("vv"),
				Type:  ssmtypes.ParameterTypeSecureString,
			},
		},
		putOut: &awsssm.PutParameterOutput{Version: 9},
	}
	s := newServiceWithAPI(fake)

	out, err := s.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Version != 9 {
		t.Errorf("Version = %d, want 9", out.Version)
	}
	if fake.putIn.Type != ssmtypes.ParameterTypeSecureString {
		t.Errorf("dest Type = %q, want SecureString (carried from source)", fake.putIn.Type)
	}
	if got := *fake.putIn.Name; got != "/dest" {
		t.Errorf("dest Name = %q, want /dest", got)
	}
}

func TestService_Copy_GetOneFailsSaveNotCalled(t *testing.T) {
	fake := &fakeSSMAPI{getErr: &ssmtypes.ParameterNotFound{}}
	s := newServiceWithAPI(fake)

	_, err := s.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if fake.putIn != nil {
		t.Error("PutParameter was called despite GetOne failure")
	}
}

func TestService_Copy_SaveFails(t *testing.T) {
	fake := &fakeSSMAPI{
		getOut: &awsssm.GetParameterOutput{
			Parameter: &ssmtypes.Parameter{
				Name: strPtr("/src"), Value: strPtr("v"), Type: ssmtypes.ParameterTypeString,
			},
		},
		putErr: &ssmtypes.ParameterLimitExceeded{},
	}
	s := newServiceWithAPI(fake)

	_, err := s.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterLimitExceeded
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterLimitExceeded", err)
	}
}
