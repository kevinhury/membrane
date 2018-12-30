package interceptors

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kevinhury/membrane/config/actions"

	"github.com/kevinhury/membrane/config"
)

// RequestModifier func
func RequestModifier(r *http.Request, pipelines []config.Pipeline) error {
	for i := range pipelines {
		plugins := pipelines[i].PluginsMatchingName("request-transform")
		for j := range plugins {
			err := modifyRequest(r, plugins[j])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func modifyRequest(r *http.Request, plugin config.Plugin) error {
	err := modifyRequestBody(r, plugin)
	if err != nil {
		return err
	}
	return nil
}

func modifyRequestBody(r *http.Request, plugin config.Plugin) error {
	action, ok := plugin.Action.(actions.RequestTransform)
	if !ok {
		return nil
	}
	if action.Body == nil {
		return nil
	}

	if action.Body.Duplicate != nil {

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		bodyFormatted := make(map[string]interface{})
		err = json.Unmarshal(b, &bodyFormatted)
		if err != nil {
			return err
		}

		for key, newKey := range action.Body.Duplicate {
			if value, found := bodyFormatted[key]; found {
				bodyFormatted[newKey] = value
			}
		}

		b, err = json.Marshal(bodyFormatted)
		if err != nil {
			return err
		}

		r.ContentLength = int64(len(b)) // TODO Check type conversion
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}

	return nil
}
