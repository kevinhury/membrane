package interceptors

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/kevinhury/membrane/config/actions"

	"github.com/kevinhury/membrane/config"
)

// RequestModifier func
func RequestModifier(r *http.Request, w http.ResponseWriter, pipelines []config.Pipeline) error {
	for i := range pipelines {
		plugins := pipelines[i].PluginsMatchingName("request-transform")
		for j := range plugins {
			err := modifyRequest(r, plugins[j])
			if err != nil {
				return err
			}
		}
		jwtPlugins := pipelines[i].PluginsMatchingName("jwt")
		for j := range jwtPlugins {
			err := jwtPlugin(r, w, jwtPlugins[j])
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
		return errors.New("")
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

func jwtPlugin(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action, ok := plugin.Action.(actions.JWT)
	if !ok {
		return errors.New("")
	}

	if len(action.Secret) == 0 || len(action.Strategy) == 0 {
		return errors.New("")
	}

	if action.Strategy == "bearer" {
		token, err := getAuthToken(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return err
		}

		err = validateToken(token, action.Secret)
		if err != nil {
			w.Write([]byte(err.Error()))
			return nil
		}
	}

	return nil
}

func getAuthToken(r *http.Request) (string, error) {
	requestToken := r.Header.Get("Authorization")
	splitToken := strings.Split(requestToken, "Bearer")
	if len(splitToken) < 2 {
		return "", errors.New("missing auth token")
	}
	return splitToken[1], nil
}

func validateToken(token, secret string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	log.Println(t)

	return err
}
