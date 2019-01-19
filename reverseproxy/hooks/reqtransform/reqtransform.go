package reqtransform

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/config/actions"
)

// Hook struct
type Hook struct {
	Config *config.Configuration
}

// PreHook func
func (h Hook) PreHook(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action := plugin.Action.(actions.RequestTransform)

	if action.Body != nil {
		err := modifyRequestBody(r, w, plugin)
		if err != nil {
			return err
		}
	}
	if action.Query != nil {
		err := modifyRequestQuery(r, w, plugin)
		if err != nil {
			return err
		}
	}

	return nil
}

func modifyRequestBody(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action := plugin.Action.(actions.RequestTransform)

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

	for key, val := range action.Body.Append {
		bodyFormatted[key] = val
	}

	b, err = json.Marshal(bodyFormatted)
	if err != nil {
		return err
	}

	r.ContentLength = int64(len(b))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return nil
}

func modifyRequestQuery(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action := plugin.Action.(actions.RequestTransform)

	q := r.URL.Query()
	for key, newKey := range action.Query.Duplicate {
		q.Add(newKey, q.Get(key))
	}

	for key, val := range action.Query.Append {
		q.Add(key, val)
	}

	r.URL.RawQuery = q.Encode()

	return nil
}
