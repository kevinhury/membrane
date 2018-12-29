package interceptors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/kevinhury/membrane/config"
)

// ResponseModifier func
func ResponseModifier(pipelines []config.Pipeline) func(*http.Response) error {
	return func(resp *http.Response) error {
		for i := range pipelines {
			plugins := pipelines[i].PluginsMatchingName("response-transform")
			for j := range plugins {
				modifyRespose(resp, plugins[j])
			}
		}
		return nil
	}
}

type pluginAction struct {
	modifyStatus int
	setHeader    map[string][]byte
}

func getAction(action interface{}) (pluginAction, bool) {
	act, ok := action.(pluginAction)
	return act, ok
}

func modifyRespose(resp *http.Response, plugin config.Plugin) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	log.Printf("Applying plugin %+v", plugin)
	log.Printf("Intercepted response %s", string(b))

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	heads := pluginHeaders(plugin.Action["setHeaders"])
	for k, v := range heads {
		resp.Header.Set(k, string(v))
	}
	if modifyStatus, ok := plugin.Action["modifyStatus"]; ok {
		header := fmt.Sprintf("%+v", modifyStatus)
		if header != "" {
			code, err := strconv.Atoi(header)
			if err == nil {
				resp.StatusCode = code
			}
		}
	}

	return nil
}

func pluginHeaders(action interface{}) map[string]string {
	mapped := action.(map[interface{}]interface{})
	headers := make(map[string]string, len(mapped))

	for k, v := range mapped {
		key := fmt.Sprintf("%+v", k)
		headers[key] = fmt.Sprintf("%+v", v)
	}

	return headers
}
