package app

import "testing"

func TestWriteConfig(t *testing.T) {
	ci := configInput{
		OrgName: "foo",
		AppName: "test-service",
		Envs: "dev,test,prod",
	}

}