package common

import "time"

// EspParam represents the parameter that is being managed by ESP
type EspParam struct {
	Id string
	Name string
	Path string
	Type string
	Secure bool
	Value string
	Version int64
	LastModifiedDate time.Time
}

// EspParamInput represents parameter that is going to be saved
type EspParamInput struct {
	Name string
	Secure bool
	Value string
}

// GetOneInput represents the query to get a param
type GetOneInput struct {
	Name string
	Decrypt bool
}

// SaveOutput represents the response from a save operation
type SaveOutput struct {
	Version int64
}

// ListParamInput represents the output of a list query
type ListParamInput struct {
	Path string
	Decrypt bool
}

type CopyCommand struct {
	Source      string
	Destination string
}
