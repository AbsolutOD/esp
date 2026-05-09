package client

import (
	"fmt"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/AbsolutOD/esp/internal/ssm"
)

// Client the main interface that is defined by the backend implementation.
type Client interface {
	Save(p common.EspParamInput) (common.SaveOutput, error)
	GetOne(p common.GetOneInput) (common.EspParam, error)
	GetMany(p common.ListParamInput) ([]common.EspParam, error)
	Copy(cc common.CopyCommand) (common.SaveOutput, error)
	Delete(p common.DeleteInput) (string, error)
}

// EspClient is the main struct for interacting with the backend driver
type EspClient struct {
	Backend string
	Client  Client
}

// New creates a new instance of the Client for esp
func New(c *app.Config) (*EspClient, error) {
	if c.Backend != "ssm" {
		return nil, fmt.Errorf("unsupported backend %q", c.Backend)
	}
	svc, err := ssm.New()
	if err != nil {
		return nil, err
	}
	return &EspClient{
		Backend: c.Backend,
		Client:  svc,
	}, nil
}

// GetParam Queries the ssm param
func (c *EspClient) GetParam(i common.GetOneInput) (common.EspParam, error) {
	in := common.GetOneInput{
		Name:    i.Name,
		Decrypt: i.Decrypt,
	}
	return c.Client.GetOne(in)
}

// ListParams takes a path and returns all of the parameters under it
func (c *EspClient) ListParams(p common.ListParamInput) ([]common.EspParam, error) {
	return c.Client.GetMany(p)
}

// Save stores the parameter in the configured backend
func (c *EspClient) Save(p common.EspParamInput) (common.SaveOutput, error) {
	return c.Client.Save(p)
}

// Delete removes a parameter from the backend
func (c *EspClient) Delete(p common.DeleteInput) (string, error) {
	return c.Client.Delete(p)
}

// Copy allows you to copy one path to another
func (c *EspClient) Copy(cc common.CopyCommand) (common.EspParam, error) {
	if _, err := c.Client.Copy(cc); err != nil {
		return common.EspParam{}, err
	}

	query := common.GetOneInput{
		Name:    cc.Destination,
		Decrypt: true,
	}
	return c.GetParam(query)
}

// Move copies a single param to a new path in the backend
func (c *EspClient) Move(mc common.MoveCommand) (common.MoveCommand, error) {
	p, err := c.GetParam(common.GetOneInput{
		Name:    mc.Source,
		Decrypt: true,
	})
	if err != nil {
		return mc, err
	}

	if _, err := c.Save(common.EspParamInput{
		Name:   mc.Destination,
		Secure: p.Secure,
		Value:  p.Value,
	}); err != nil {
		return mc, err
	}

	if _, err := c.Delete(common.DeleteInput{Name: mc.Source}); err != nil {
		return mc, err
	}
	return mc, nil
}
