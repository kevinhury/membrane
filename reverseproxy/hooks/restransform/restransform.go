package restransform

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/config/actions"
)

// Hook struct
type Hook struct {
	Config *config.Configuration
}

// PostHook func
func (h Hook) PostHook(resp *http.Response, plugin config.Plugin) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	action := plugin.Action.(actions.ResponseTransform)
	for k, v := range action.SetHeaders {
		resp.Header.Set(k, string(v))
	}

	if action.ModifyStatus != 0 {
		resp.StatusCode = action.ModifyStatus
	}

	if ctype := resp.Header.Get("content-type"); !strings.Contains(ctype, "application/json") {
		log.Printf("Could not reformat body. got %s", ctype)
	} else {
		bodyJSON := make(map[string]interface{})
		err := json.Unmarshal(b, &bodyJSON)
		if err == nil {
			for k, v := range action.ReformatBody {
				if value, ok := bodyJSON[k]; ok {
					delete(bodyJSON, k)
					bodyJSON[v] = value
				}
			}
		}
		convertedB, err := json.Marshal(bodyJSON)
		if err == nil {
			b = convertedB
		}
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	return nil
}
