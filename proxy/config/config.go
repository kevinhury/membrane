package config

import (
	"github.com/go-yaml/yaml"
)

// Configuration struct
type Configuration struct {
	APIEndpoints []struct {
		Name  string `yaml:"name"`
		Host  string `yaml:"host"`
		Paths string `yaml:"paths"`
	} `yaml:"apiEndpoints"`
	ServiceEndpoints []struct {
		Name string `yaml:"name"`
		Uerl string `yaml:"uerl"`
	} `yaml:"serviceEndpoints"`
	Pipelines []struct {
		Name         string   `yaml:"name"`
		APIEndpoints []string `yaml:"apiEndpoints"`
		Policies     []struct {
			Proxy []struct {
				Action struct {
					ServiceEndpoint string `yaml:"serviceEndpoint"`
					ChangeOrigin    bool   `yaml:"changeOrigin"`
				} `yaml:"action"`
			} `yaml:"proxy"`
		} `yaml:"policies"`
	} `yaml:"pipelines"`
}

// // Configuration struct
// type Configuration struct {
// 	APIVersion string `yaml:"apiVersion"`
// 	Kind       string `yaml:"kind"`
// 	Metadata   struct {
// 		Name      string `yaml:"name"`
// 		Namespace string `yaml:"namespace"`
// 		Labels    struct {
// 			RouterDeisIoRoutable string `yaml:"router.deis.io/routable"`
// 		} `yaml:"labels"`
// 		Annotations struct {
// 			RouterDeisIoDomains string `yaml:"router.deis.io/domains"`
// 		} `yaml:"annotations"`
// 	} `yaml:"metadata"`
// 	Spec struct {
// 		Type     string `yaml:"type"`
// 		Selector struct {
// 			App string `yaml:"app"`
// 		} `yaml:"selector"`
// 		Ports []struct {
// 			Name       string `yaml:"name"`
// 			Port       int    `yaml:"port"`
// 			TargetPort int    `yaml:"targetPort"`
// 			NodePort   int    `yaml:"nodePort,omitempty"`
// 		} `yaml:"ports"`
// 	} `yaml:"spec"`
// }
// `
// apiVersion: v1
// kind: Service
// metadata:
//   name: myName
//   namespace: default
//   labels:
//     router.deis.io/routable: "true"
//   annotations:
//     router.deis.io/domains: ""
// spec:
//   type: NodePort
//   selector:
//     app: myName
//   ports:
//     - name: http
//       port: 80
//       targetPort: 80
//     - name: https
//       port: 443
//       targetPort: 443
// `

// Parse yaml to Configuration
func Parse(data []byte) (*Configuration, error) {
	var conf Configuration
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
