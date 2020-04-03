package client

import (
	"github.com/absolutod/esp/internal/common"
	"github.com/absolutod/esp/internal/ssm"
)

type Backend string

type Client interface {
	Save(p common.EspParam) (common.EspParam, error)
	GetOne(p common.GetOneInput) (common.EspParam, error)
}

type EspClient struct {
	Backend Backend
	Client Client
}

func New(c EspClient) *EspClient {
	if c.Backend == "ssm" {
		svc := ssm.New()
		c.Client = svc
	}
	return &c
}

// getParam Queries the ssm param
func (c *EspClient) GetParam(debug bool, key string) common.EspParam {
	in := common.GetOneInput{
		Name: key,
		Decrypt: debug,
	}
	param, err := c.Client.GetOne(in)
	if err != nil {
		panic(err)
		//fmt.Println(err)
	}

	return param
}

