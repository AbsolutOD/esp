package client

import (
	"fmt"
	"github.com/absolutod/esp/internal/errors"
	"github.com/pinpt/esp/internal/ssm"
)

type Backend string

type Client interface {
	Save(p EspParam) (EspParam, error)
	GetOne(p GetOneInput) (EspParam, error)
}

type EspParam struct {
	Id string
	Name string
	Path string
	Type string
	Secure bool
	Value string
	Version int
	LastModifiedDate float32
}

type EspParamInput struct {
	Path string
	Name string
	Secure bool
	Value string
}

type GetOneInput struct {
	Name string
	Decrypt bool
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

	c.Client.Init()
	return &c
}

// getParam Queries the ssm param
func (c *EspClient) GetParam(debug bool, key string) EspParam {
	in := GetOneInput{
		Name: key,
		Decrypt: debug,
	}
	param, err := c.Client.GetOne(in)
	if err != nil {
		fmt.Println(err)
	}

	return param
}

