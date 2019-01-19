package jwt

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/config/actions"
)

// Hook struct
type Hook struct {
	Config *config.Configuration
}

// PreHook func
func (h Hook) PreHook(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action := plugin.Action.(actions.JWT)

	if action.Secret == "" || action.Strategy == "" {
		return errors.New("JWT Plugin: Bad configuration")
	}

	if action.Strategy == "bearer" {
		token, err := getAuthToken(r)
		if err != nil {
			writeJSONError(w, err)
			return err
		}

		err = validateToken(token, action.Secret)
		if err != nil {
			writeJSONError(w, err)
			return err
		}
	}

	return nil
}

func getAuthToken(r *http.Request) (string, error) {
	requestToken := r.Header.Get("Authorization")
	splitToken := strings.Split(requestToken, "Bearer ")
	if len(splitToken) < 2 {
		return "", errors.New("missing auth token")
	}
	return splitToken[1], nil
}

func validateToken(token, secret string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return err
}

func writeJSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	defaultErrorMsg := map[string]string{
		"error": err.Error(),
	}
	bs, _ := json.Marshal(defaultErrorMsg)

	w.Write(bs)
}
