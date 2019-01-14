package jwt

import (
	"errors"
	"log"
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
