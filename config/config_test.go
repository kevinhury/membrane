package config

import (
	"testing"
)

func TestParsing(t *testing.T) {
	mockConfig := `
inboundEndpoints:
  - name: api
    host: 'localhost'
    paths: '/ip'

outboundEndpoints:
  - name: httpbin
    uerl: 'https://httpbin.org' 

pipelines:
    - name: getting-started
      inboundEndpoints:
        - api
      policies:
        - proxy:
            - action:
                outboundEndpoint: httpbin
                changeOrigin: true

  `
	conf, err := Parse([]byte(mockConfig))
	if err != nil {
		t.Errorf("Expected to parse, got error %+v", err)
	}
	if conf == nil {
		t.Error("Expected conf got nil")
	}
}
