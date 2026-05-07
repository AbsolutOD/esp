package ssm

import (
	"errors"
	"testing"

	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
)

// genericAPIErr returns a *smithy.GenericAPIError (which satisfies
// smithy.APIError) for use in fallthrough cases.
func genericAPIErr() error {
	return &smithy.GenericAPIError{Code: "Unknown", Message: "unknown api failure"}
}

// plainErr returns a non-API error to verify the no-match-returns-nil path.
func plainErr() error { return errors.New("plain non-api error") }

func TestCheckBaseSSMErrors(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantSelf  bool // expect the function to return err unchanged
		wantNil   bool // expect nil
	}{
		{name: "InvalidKeyId is mapped", err: &ssmtypes.InvalidKeyId{}, wantSelf: true},
		{name: "InternalServerError is mapped", err: &ssmtypes.InternalServerError{}, wantSelf: true},
		{name: "generic smithy.APIError falls through and is returned", err: genericAPIErr(), wantSelf: true},
		{name: "plain non-API error returns nil", err: plainErr(), wantNil: true},
		{name: "nil input returns nil", err: nil, wantNil: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkBaseSSMErrors(tc.err)
			if tc.wantSelf && got != tc.err {
				t.Errorf("checkBaseSSMErrors(%v) = %v, want input err returned", tc.err, got)
			}
			if tc.wantNil && got != nil {
				t.Errorf("checkBaseSSMErrors(%v) = %v, want nil", tc.err, got)
			}
		})
	}
}

func TestCheckSSMGetParameterError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantSelf bool
		wantNil  bool
	}{
		{name: "ParameterNotFound is mapped", err: &ssmtypes.ParameterNotFound{}, wantSelf: true},
		{name: "generic smithy.APIError falls through and is returned", err: genericAPIErr(), wantSelf: true},
		{name: "plain non-API error returns nil", err: plainErr(), wantNil: true},
		{name: "nil input returns nil", err: nil, wantNil: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkSSMGetParameterError(tc.err)
			if tc.wantSelf && got != tc.err {
				t.Errorf("checkSSMGetParameterError(%v) = %v, want input err returned", tc.err, got)
			}
			if tc.wantNil && got != nil {
				t.Errorf("checkSSMGetParameterError(%v) = %v, want nil", tc.err, got)
			}
		})
	}
}

func TestCheckSSMPutParameterError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantSelf bool
		wantNil  bool
	}{
		{name: "ParameterLimitExceeded", err: &ssmtypes.ParameterLimitExceeded{}, wantSelf: true},
		{name: "TooManyUpdates", err: &ssmtypes.TooManyUpdates{}, wantSelf: true},
		{name: "HierarchyTypeMismatchException", err: &ssmtypes.HierarchyTypeMismatchException{}, wantSelf: true},
		{name: "InvalidAllowedPatternException", err: &ssmtypes.InvalidAllowedPatternException{}, wantSelf: true},
		{name: "ParameterMaxVersionLimitExceeded", err: &ssmtypes.ParameterMaxVersionLimitExceeded{}, wantSelf: true},
		{name: "UnsupportedParameterType", err: &ssmtypes.UnsupportedParameterType{}, wantSelf: true},
		{name: "PoliciesLimitExceededException", err: &ssmtypes.PoliciesLimitExceededException{}, wantSelf: true},
		{name: "InvalidPolicyTypeException", err: &ssmtypes.InvalidPolicyTypeException{}, wantSelf: true},
		{name: "InvalidPolicyAttributeException", err: &ssmtypes.InvalidPolicyAttributeException{}, wantSelf: true},
		{name: "IncompatiblePolicyException", err: &ssmtypes.IncompatiblePolicyException{}, wantSelf: true},
		{name: "ParameterAlreadyExists", err: &ssmtypes.ParameterAlreadyExists{}, wantSelf: true},
		{name: "generic smithy.APIError falls through", err: genericAPIErr(), wantSelf: true},
		{name: "plain non-API error returns nil", err: plainErr(), wantNil: true},
		{name: "nil input returns nil", err: nil, wantNil: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkSSMPutParameterError(tc.err)
			if tc.wantSelf && got != tc.err {
				t.Errorf("checkSSMPutParameterError(%v) = %v, want input err returned", tc.err, got)
			}
			if tc.wantNil && got != nil {
				t.Errorf("checkSSMPutParameterError(%v) = %v, want nil", tc.err, got)
			}
		})
	}
}

func TestCheckDeleteParameterError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantSelf bool
		wantNil  bool
	}{
		{name: "ParameterNotFound is mapped", err: &ssmtypes.ParameterNotFound{}, wantSelf: true},
		{name: "generic smithy.APIError falls through", err: genericAPIErr(), wantSelf: true},
		{name: "plain non-API error returns nil", err: plainErr(), wantNil: true},
		{name: "nil input returns nil", err: nil, wantNil: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkDeleteParameterError(tc.err)
			if tc.wantSelf && got != tc.err {
				t.Errorf("checkDeleteParameterError(%v) = %v, want input err returned", tc.err, got)
			}
			if tc.wantNil && got != nil {
				t.Errorf("checkDeleteParameterError(%v) = %v, want nil", tc.err, got)
			}
		})
	}
}

// checkSSMByPathError differs from the others: it does NOT have a
// trailing smithy.APIError fallback, so a generic API error returns
// nil rather than passing through. The test pins this asymmetry.
func TestCheckSSMByPathError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantSelf bool
		wantNil  bool
	}{
		{name: "InternalServerError", err: &ssmtypes.InternalServerError{}, wantSelf: true},
		{name: "InvalidFilterKey", err: &ssmtypes.InvalidFilterKey{}, wantSelf: true},
		{name: "InvalidFilterOption", err: &ssmtypes.InvalidFilterOption{}, wantSelf: true},
		{name: "InvalidFilterValue", err: &ssmtypes.InvalidFilterValue{}, wantSelf: true},
		{name: "InvalidKeyId", err: &ssmtypes.InvalidKeyId{}, wantSelf: true},
		{name: "InvalidNextToken", err: &ssmtypes.InvalidNextToken{}, wantSelf: true},
		{name: "generic smithy.APIError returns nil (no fallback in this function)", err: genericAPIErr(), wantNil: true},
		{name: "plain non-API error returns nil", err: plainErr(), wantNil: true},
		{name: "nil input returns nil", err: nil, wantNil: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkSSMByPathError(tc.err)
			if tc.wantSelf && got != tc.err {
				t.Errorf("checkSSMByPathError(%v) = %v, want input err returned", tc.err, got)
			}
			if tc.wantNil && got != nil {
				t.Errorf("checkSSMByPathError(%v) = %v, want nil", tc.err, got)
			}
		})
	}
}

// TestCheckSSMError exercises the dispatcher: each known action
// routes to the matching per-action function; an unknown action
// (e.g. GetMany, which has no case) falls through to checkBaseSSMErrors.
func TestCheckSSMError(t *testing.T) {
	tests := []struct {
		name     string
		action   action
		err      error
		wantSelf bool
		wantNil  bool
	}{
		{
			name:     "Get routes to checkSSMGetParameterError",
			action:   Get,
			err:      &ssmtypes.ParameterNotFound{},
			wantSelf: true,
		},
		{
			name:     "Save routes to checkSSMPutParameterError",
			action:   Save,
			err:      &ssmtypes.ParameterAlreadyExists{},
			wantSelf: true,
		},
		{
			name:     "Delete routes to checkDeleteParameterError",
			action:   Delete,
			err:      &ssmtypes.ParameterNotFound{},
			wantSelf: true,
		},
		{
			name:     "GetMany action falls through to base error mapper",
			action:   GetMany,
			err:      &ssmtypes.InternalServerError{},
			wantSelf: true,
		},
		{
			name:    "GetMany action with non-API error returns nil",
			action:  GetMany,
			err:     plainErr(),
			wantNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkSSMError(tc.action, tc.err)
			if tc.wantSelf && got != tc.err {
				t.Errorf("checkSSMError(%q, %v) = %v, want input err returned", tc.action, tc.err, got)
			}
			if tc.wantNil && got != nil {
				t.Errorf("checkSSMError(%q, %v) = %v, want nil", tc.action, tc.err, got)
			}
		})
	}
}
