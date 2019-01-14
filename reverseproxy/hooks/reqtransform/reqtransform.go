package reqtransform

import (
	"bytes"
	"encoding/json"
	"errors"
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
	action, ok := plugin.Action.(actions.RequestTransform)
	if !ok {
		return errors.New("Unsupported action")
	}
	if action.Body != nil {
		err := modifyRequestBody(r, w, plugin)
		if err != nil {
			return err
		}
	}
	if action.Query != nil {
		return nil
	}

	return nil
}

func modifyRequestBody(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action := plugin.Action.(actions.RequestTransform)

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

		r.ContentLength = int64(len(b))
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}

	return nil
}
