package app

import "testing"

func TestWriteConfig(t *testing.T) {
	c := New(false )
	ci := configInput{
		Backend: "ssm",
		OrgName: "foo",
		AppName: "test-service",
		Envs: "dev,test,prod",
	}

	c.WriteConfig(&ci)

}