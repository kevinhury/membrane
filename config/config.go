package config

import (
	"os"

	"github.com/go-yaml/yaml"
	"github.com/kevinhury/membrane/config/actions"
	"github.com/kevinhury/membrane/config/urlutils"
	"github.com/mitchellh/mapstructure"
)

// Configuration struct
type Configuration struct {
	ConfigMap *Map
	endpoints map[string]*InboundEndpoint
	services  map[string]*OutboundEndpoint
}

// NewWithData func
func NewWithData(data []byte) (*Configuration, error) {
	cm, err := Parse(data)
	if err != nil {
		return nil, err
	}

	config := &Configuration{ConfigMap: cm}
	config.endpoints = make(map[string]*InboundEndpoint, len(config.ConfigMap.InboundEndpoints))
	for i := 0; i < len(config.ConfigMap.InboundEndpoints); i++ {
		ep := config.ConfigMap.InboundEndpoints[i]
		config.endpoints[ep.Name] = &ep
	}
	config.services = make(map[string]*OutboundEndpoint, len(config.ConfigMap.OutboundEndpoints))
	for i := 0; i < len(config.ConfigMap.OutboundEndpoints); i++ {
		out := config.ConfigMap.OutboundEndpoints[i]
		config.services[out.Name] = &out
	}

	return config, nil
}

// Pipelines func
func (c *Configuration) Pipelines(host, path, method string) []Pipeline {
	var pipelines []Pipeline
	endpoints := c.Endpoints(host, path, method)

	for _, p := range c.ConfigMap.Pipelines {
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
	for _, ep := range c.ConfigMap.InboundEndpoints {
		if ep.Host != "" {
			if host != ep.Host {
				continue
			}
		}

		if len(ep.Paths) > 0 {
			_, found := urlutils.MatchPath(path, ep.Paths)
			if found == false {
				continue
			}
		}

		if len(ep.Methods) > 0 {
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

// DefaultHost func
func (ie *InboundEndpoint) DefaultHost() string {
	name, _ := os.Hostname()
	return name
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

// Plugin func
func (p *Pipeline) Plugin(name string) *Plugin {
	for _, plugin := range p.Plugins {
		if plugin.Name == name {
			return &plugin
		}
	}
	return nil
}

// PluginsMatchingName func
func (p *Pipeline) PluginsMatchingName(name string) []Plugin {
	var plugins []Plugin
	for _, p := range p.Plugins {
		if p.Name == name {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// Plugin struct
type Plugin struct {
	Name       string                 `yaml:"name"`
	Conditions map[string]interface{} `yaml:"conditions"`
	Action     interface{}            `yaml:"action"`
}

// Parse yaml to Configuration
func Parse(data []byte) (*Map, error) {
	var conf Map
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(conf.Pipelines); i++ {
		pipeline := conf.Pipelines[i]
		for idx := 0; idx < len(pipeline.Plugins); idx++ {
			plugin := pipeline.Plugins[idx]
			if plugin.Name == "jwt" {
				var act actions.JWT
				mapstructure.Decode(plugin.Action, &act)
				pipeline.Plugins[idx].Action = act
			} else if plugin.Name == "jwt-extract" {
				var act actions.JWTExtract
				mapstructure.Decode(plugin.Action, &act)
				pipeline.Plugins[idx].Action = act
			} else if plugin.Name == "proxy" {
				var act actions.Proxy
				mapstructure.Decode(plugin.Action, &act)
				pipeline.Plugins[idx].Action = act
			} else if plugin.Name == "response-transform" {
				var act actions.ResponseTransform
				mapstructure.Decode(plugin.Action, &act)
				pipeline.Plugins[idx].Action = act
			} else if plugin.Name == "request-transform" {
				var act actions.RequestTransform
				mapstructure.Decode(plugin.Action, &act)
				pipeline.Plugins[idx].Action = act
			} else if plugin.Name == "cors" {
				var act actions.Cors
				mapstructure.Decode(plugin.Action, &act)
				pipeline.Plugins[idx].Action = act
			}
		}
	}

	return &conf, nil
}
