package ssm

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/AbsolutOD/esp/internal/common"
)

func TestSelectType(t *testing.T) {
	if got := selectType(true); got != ssmtypes.ParameterTypeSecureString {
		t.Errorf("selectType(true) = %v, want %v", got, ssmtypes.ParameterTypeSecureString)
	}
	if got := selectType(false); got != ssmtypes.ParameterTypeString {
		t.Errorf("selectType(false) = %v, want %v", got, ssmtypes.ParameterTypeString)
	}
}

func TestConvertToEspParam(t *testing.T) {
	modified := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		in   ssmtypes.Parameter
		want common.EspParam
	}{
		{
			name: "secure string sets Secure: true",
			in: ssmtypes.Parameter{
				ARN:              aws.String("arn:aws:ssm:us-east-1:1:parameter/acme/dev/billing/SECRET"),
				Name:             aws.String("/acme/dev/billing/SECRET"),
				Type:             ssmtypes.ParameterTypeSecureString,
				Value:            aws.String("hunter2"),
				Version:          7,
				LastModifiedDate: aws.Time(modified),
			},
			want: common.EspParam{
				Id:               "arn:aws:ssm:us-east-1:1:parameter/acme/dev/billing/SECRET",
				Name:             "/acme/dev/billing/SECRET",
				Type:             "SecureString",
				Value:            "hunter2",
				Version:          7,
				LastModifiedDate: modified,
				Secure:           true,
			},
		},
		{
			name: "plain string leaves Secure: false",
			in: ssmtypes.Parameter{
				ARN:              aws.String("arn:plain"),
				Name:             aws.String("DB_URL"),
				Type:             ssmtypes.ParameterTypeString,
				Value:            aws.String("postgres://"),
				Version:          1,
				LastModifiedDate: aws.Time(modified),
			},
			want: common.EspParam{
				Id:               "arn:plain",
				Name:             "DB_URL",
				Type:             "String",
				Value:            "postgres://",
				Version:          1,
				LastModifiedDate: modified,
				Secure:           false,
			},
		},
		{
			name: "stringlist leaves Secure: false (only SecureString flips it)",
			in: ssmtypes.Parameter{
				ARN:   aws.String("arn:list"),
				Name:  aws.String("LIST"),
				Type:  ssmtypes.ParameterTypeStringList,
				Value: aws.String("a,b,c"),
			},
			want: common.EspParam{
				Id:     "arn:list",
				Name:   "LIST",
				Type:   "StringList",
				Value:  "a,b,c",
				Secure: false,
			},
		},
		{
			name: "nil pointer fields render as zero values via aws.ToString / aws.ToTime",
			in: ssmtypes.Parameter{
				Type: ssmtypes.ParameterTypeString,
			},
			want: common.EspParam{
				Type: "String",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := convertToEspParam(tc.in)
			if got != tc.want {
				t.Errorf("convertToEspParam(%+v) =\n  %+v\nwant\n  %+v", tc.in, got, tc.want)
			}
		})
	}
}
