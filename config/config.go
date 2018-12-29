package config

import (
	"github.com/go-yaml/yaml"
)

// Configuration struct
type Configuration struct {
	configMap *Map
	endpoints map[string]*InboundEndpoint
	services  map[string]*OutboundEndpoint
}

// NewWithData func
func NewWithData(data []byte) (*Configuration, error) {
	cm, err := Parse(data)
	if err != nil {
		return nil, err
	}

	config := &Configuration{configMap: cm}
	config.endpoints = make(map[string]*InboundEndpoint, len(config.configMap.InboundEndpoints))
	for _, ep := range config.configMap.InboundEndpoints {
		config.endpoints[ep.Name] = &ep
	}
	config.services = make(map[string]*OutboundEndpoint, len(config.configMap.OutboundEndpoints))
	for _, se := range config.configMap.OutboundEndpoints {
		config.services[se.Name] = &se
	}

	return config, nil
}

// Pipelines func
func (c *Configuration) Pipelines(host, path, method string) []Pipeline {
	var pipelines []Pipeline
	endpoints := c.Endpoints(host, path, method)

	for _, p := range c.configMap.Pipelines {
		for _, epName := range p.InboundEndpoints {
			if _, ok := endpoints[epName]; ok {
				pipelines = append(pipelines, p)
			}
		}
	}

	return pipelines
}

// Endpoints func
func (c *Configuration) Endpoints(host, path, method string) map[string]*InboundEndpoint {
	endpoints := make(map[string]*InboundEndpoint, 0)
	for _, ep := range c.configMap.InboundEndpoints {
		if host != ep.Host {
			continue
		}
		matchP := false
		for _, p := range ep.Paths {
			if path == p {
				matchP = true
				break
			}
		}
		if matchP != true {
			continue
		}
		matchM := false
		for _, m := range ep.Methods {
			if method == m {
				matchM = true
				break
			}
		}
		if matchM != true {
			continue
		}
		endpoints[ep.Name] = &ep
	}

	return endpoints
}

// Service func
func (c *Configuration) Service(name string) *OutboundEndpoint {
	return c.services[name]
}

// Map struct
type Map struct {
	InboundEndpoints  []InboundEndpoint  `yaml:"inboundEndpoints"`
	OutboundEndpoints []OutboundEndpoint `yaml:"outboundEndpoints"`
	Pipelines         []Pipeline         `yaml:"pipelines"`
}

// InboundEndpoint struct
type InboundEndpoint struct {
	Name    string   `yaml:"name"`
	Host    string   `yaml:"host"`
	Paths   []string `yaml:"paths"`
	Methods []string `yaml:"methods"`
}

// OutboundEndpoint struct
type OutboundEndpoint struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Pipeline struct
type Pipeline struct {
	Name             string   `yaml:"name"`
	InboundEndpoints []string `yaml:"inboundEndpoints"`
	Plugins          []Plugin `yaml:"plugins"`
}

// Plugin struct
type Plugin struct {
	Name       string                 `yaml:"name"`
	Conditions map[string]interface{} `yaml:"conditions"`
	Action     map[string]interface{} `yaml:"action"`
}

// Parse yaml to Configuration
func Parse(data []byte) (*Map, error) {
	var conf Map
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	// TODO: VALIDATE CONFIG

	return &conf, nil
}
