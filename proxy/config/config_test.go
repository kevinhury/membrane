package config

import (
	"testing"
)

func TestParsing(t *testing.T) {
	mockConfig := `
apiEndpoints:
  - name: api
    host: 'localhost'
    paths: '/ip'

serviceEndpoints:
  - name: httpbin
    uerl: 'https://httpbin.org' 

pipelines:
    - name: getting-started
      apiEndpoints:
        - api
      policies:
        - proxy:
            - action:
                serviceEndpoint: httpbin
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
