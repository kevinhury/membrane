package actions

import (
	"github.com/mitchellh/mapstructure"
)

// JWT struct
type JWT struct {
	Secret   string `yaml:"secret"`
	Strategy string `yaml:"strategy"`
}

// JWTExtract struct
type JWTExtract struct {
	Secret string            `yaml:"secret"`
	Body   map[string]string `yaml:"body"`
	Query  map[string]string `yaml:"query"`
}

// Proxy struct
type Proxy struct {
	OutboundEndpoint string `yaml:"outboundEndpoint"`
	KeepOrigin       bool   `yaml:"keepOrigin"`
	prependPath      bool   `yaml:"prependPath"`
}

// ResponseTransform struct
type ResponseTransform struct {
	ModifyStatus int               `yaml:"modifyStatus"`
	SetHeaders   map[string]string `yaml:"setHeaders"`
	ReformatBody map[string]string `yaml:"reformatBody"`
}

// Transform struct
type Transform struct {
	Append    map[string]string `yaml:"append"`
	Duplicate map[string]string `yaml:"duplicate"`
}

// RequestTransform struct
type RequestTransform struct {
	Body  *Transform `yaml:"body"`
	Query *Transform `yaml:"query"`
}

// RateLimit struct
type RateLimit struct {
	Max      int `yaml:"max"`
	WindowMs int `yaml:"windowMs"`
}

// Cors struct
type Cors struct {
	Origin  string `yaml:"origin"`
	Methods string `yaml:"methods"`
	Headers string `yaml:"headers"`
}

// Parse func
func Parse(name string, action interface{}) interface{} {
	if name == "jwt" {
		var act JWT
		mapstructure.Decode(action, &act)
		return act
	} else if name == "jwt-extract" {
		var act JWTExtract
		mapstructure.Decode(action, &act)
		return act
	} else if name == "proxy" {
		var act Proxy
		mapstructure.Decode(action, &act)
		return act
	} else if name == "response-transform" {
		var act ResponseTransform
		mapstructure.Decode(action, &act)
		return act
	} else if name == "request-transform" {
		var act RequestTransform
		mapstructure.Decode(action, &act)
		return act
	} else if name == "cors" {
		var act Cors
		mapstructure.Decode(action, &act)
		return act
	}

	return nil
}
