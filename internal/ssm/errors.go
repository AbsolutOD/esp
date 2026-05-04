package ssm

import (
	"errors"

	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
)

// checkSSMError is the entry point for check all of the based and call specific errors
func checkSSMError(a action, err error) error {
	awsErr := checkBaseSSMErrors(err)
	switch a {
	case Get:
		return checkSSMGetParameterError(err)
	case Save:
		return checkSSMPutParameterError(err)
	case Delete:
		return checkDeleteParameterError(err)
	}
	return awsErr
}

// checkBaseSSMErrors checks for the common errors all SSM API calls might return
func checkBaseSSMErrors(err error) error {
	var invalidKey *ssmtypes.InvalidKeyId
	if errors.As(err, &invalidKey) {
		return err
	}
	var internalErr *ssmtypes.InternalServerError
	if errors.As(err, &internalErr) {
		return err
	}
	// Generic API error fallback (covers MissingRegion-style errors and any unmapped types).
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

func checkDeleteParameterError(err error) error {
	var notFound *ssmtypes.ParameterNotFound
	if errors.As(err, &notFound) {
		return err
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

// checkSSMGetParameterError checks for errors the GetParameter API call might return
func checkSSMGetParameterError(err error) error {
	var notFound *ssmtypes.ParameterNotFound
	if errors.As(err, &notFound) {
		return err
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

// checkSSMPutParameterError checks for errors the PutParameter API call might return
func checkSSMPutParameterError(err error) error {
	var limitExceeded *ssmtypes.ParameterLimitExceeded
	if errors.As(err, &limitExceeded) {
		return err
	}
	var tooMany *ssmtypes.TooManyUpdates
	if errors.As(err, &tooMany) {
		return err
	}
	var hierarchyMismatch *ssmtypes.HierarchyTypeMismatchException
	if errors.As(err, &hierarchyMismatch) {
		return err
	}
	var invalidPattern *ssmtypes.InvalidAllowedPatternException
	if errors.As(err, &invalidPattern) {
		return err
	}
	var maxVersion *ssmtypes.ParameterMaxVersionLimitExceeded
	if errors.As(err, &maxVersion) {
		return err
	}
	var unsupportedType *ssmtypes.UnsupportedParameterType
	if errors.As(err, &unsupportedType) {
		return err
	}
	var policyLimit *ssmtypes.PoliciesLimitExceededException
	if errors.As(err, &policyLimit) {
		return err
	}
	var invalidPolicyType *ssmtypes.InvalidPolicyTypeException
	if errors.As(err, &invalidPolicyType) {
		return err
	}
	var invalidPolicyAttr *ssmtypes.InvalidPolicyAttributeException
	if errors.As(err, &invalidPolicyAttr) {
		return err
	}
	var incompatible *ssmtypes.IncompatiblePolicyException
	if errors.As(err, &incompatible) {
		return err
	}
	var alreadyExists *ssmtypes.ParameterAlreadyExists
	if errors.As(err, &alreadyExists) {
		return err
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

// checkSSMByPathError checks for errors the GetParameterByPath API call might return
func checkSSMByPathError(err error) error {
	var internalErr *ssmtypes.InternalServerError
	if errors.As(err, &internalErr) {
		return err
	}
	var invalidFilterKey *ssmtypes.InvalidFilterKey
	if errors.As(err, &invalidFilterKey) {
		return err
	}
	var invalidFilterOption *ssmtypes.InvalidFilterOption
	if errors.As(err, &invalidFilterOption) {
		return err
	}
	var invalidFilterValue *ssmtypes.InvalidFilterValue
	if errors.As(err, &invalidFilterValue) {
		return err
	}
	var invalidKey *ssmtypes.InvalidKeyId
	if errors.As(err, &invalidKey) {
		return err
	}
	var invalidNextToken *ssmtypes.InvalidNextToken
	if errors.As(err, &invalidNextToken) {
		return err
	}
	return nil
}
