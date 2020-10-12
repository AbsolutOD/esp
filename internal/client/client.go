package client

import (
	"github.com/pinpt/esp/internal/app"
	"github.com/pinpt/esp/internal/common"
	"github.com/pinpt/esp/internal/ssm"
)

// Client the main interface that is defined by the backend implementation.
type Client interface {
	Save(p common.EspParamInput) common.SaveOutput
	GetOne(p common.GetOneInput) common.EspParam
	GetMany(p common.ListParamInput) []common.EspParam
	Copy(cc common.CopyCommand) common.SaveOutput
	Delete(p common.DeleteInput) string
}

// EspClient is the main struct for interacting with the backend driver
type EspClient struct {
	Backend string
	Client  Client
}

// New creates a new instance of the Client for esp
func New(c *app.Config) *EspClient {
	if c.Backend == "ssm" {
		svc := ssm.New()
		svc.Init()
		return &EspClient{
			Backend: c.Backend,
			Client: svc,
		}
	} else {
		panic("Currently only the ssm backend is valid.")
	}
}

// GetParam Queries the ssm param
func (c *EspClient) GetParam(i common.GetOneInput) common.EspParam {
	in := common.GetOneInput{
		Name:    i.Name,
		Decrypt: i.Decrypt,
	}
	return c.Client.GetOne(in)
}

// ListParams takes a path and returns all of the parameters under it
func (c *EspClient) ListParams(p common.ListParamInput) []common.EspParam {
	return c.Client.GetMany(p)
}

// Save stores the parameter in the configured backend
func (c *EspClient) Save(p common.EspParamInput) common.SaveOutput {
	return c.Client.Save(p)
}

// Delete removes a parameter from the backend
func (c *EspClient) Delete(p common.DeleteInput) string {
	return c.Client.Delete(p)
}

// Copy allows you to copy one path to another
func (c *EspClient) Copy(cc common.CopyCommand) common.EspParam {
	_ = c.Client.Copy(cc)

	query := common.GetOneInput{
		Name:    cc.Destination,
		Decrypt: true,
	}
	return c.GetParam(query)
}

// Move copies a single param to a new path in the backend
func (c *EspClient) Move(mc common.MoveCommand) common.MoveCommand {
	p := c.GetParam(common.GetOneInput{
		Name:    mc.Source,
		Decrypt: true,
	})

	_ = c.Save(common.EspParamInput{
		Name:   mc.Destination,
		Secure: p.Secure,
		Value:  p.Value,
	})

	_ = c.Delete(common.DeleteInput{Name: mc.Source})
	return mc
}
