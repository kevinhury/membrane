package config

import (
	"github.com/go-yaml/yaml"
)

// Configuration struct
type Configuration struct {
	configMap *Map
	endpoints map[string]*APIEndpoint
	services  map[string]*ServiceEndpoint
}

// NewWithData func
func NewWithData(data []byte) (*Configuration, error) {
	cm, err := Parse(data)
	if err != nil {
		return nil, err
	}

	config := &Configuration{configMap: cm}
	config.endpoints = make(map[string]*APIEndpoint, len(config.configMap.APIEndpoints))
	for _, ep := range config.configMap.APIEndpoints {
		config.endpoints[ep.Name] = &ep
	}
	config.services = make(map[string]*ServiceEndpoint, len(config.configMap.ServiceEndpoints))
	for _, se := range config.configMap.ServiceEndpoints {
		config.services[se.Name] = &se
	}

	return config, nil
}

// Pipelines func
func (c *Configuration) Pipelines(host, path, method string) []Pipeline {
	var pipelines []Pipeline
	endpoints := c.Endpoints(host, path, method)

	for _, p := range c.configMap.Pipelines {
		for _, epName := range p.APIEndpoints {
			if _, ok := endpoints[epName]; ok {
				pipelines = append(pipelines, p)
			}
		}
	}

	return pipelines
}

// Endpoints func
func (c *Configuration) Endpoints(host, path, method string) map[string]*APIEndpoint {
	endpoints := make(map[string]*APIEndpoint, 0)
	for _, ep := range c.configMap.APIEndpoints {
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
func (c *Configuration) Service(name string) *ServiceEndpoint {
	return c.services[name]
}

// Map struct
type Map struct {
	APIEndpoints     []APIEndpoint     `yaml:"apiEndpoints"`
	ServiceEndpoints []ServiceEndpoint `yaml:"serviceEndpoints"`
	Pipelines        []Pipeline        `yaml:"pipelines"`
}

// APIEndpoint struct
type APIEndpoint struct {
	Name    string   `yaml:"name"`
	Host    string   `yaml:"host"`
	Paths   []string `yaml:"paths"`
	Methods []string `yaml:"methods"`
}

// ServiceEndpoint struct
type ServiceEndpoint struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Pipeline struct
type Pipeline struct {
	Name         string   `yaml:"name"`
	APIEndpoints []string `yaml:"apiEndpoints"`
	Policies     []struct {
		Name   string `yaml:"name"`
		Action struct {
			ServiceEndpoint string `yaml:"serviceEndpoint"`
			ChangeOrigin    bool   `yaml:"changeOrigin"`
		} `yaml:"action"`
	} `yaml:"policies"`
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
