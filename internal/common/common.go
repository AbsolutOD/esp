package common

import "time"

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

