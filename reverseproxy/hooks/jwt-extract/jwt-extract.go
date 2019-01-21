package jwt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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

// PreHook extracts jwt claims and sets the request body and query
func (h Hook) PreHook(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action := plugin.Action.(actions.JWTExtract)

	claims, err := getJWTClaims(r, &action)
	if err != nil {
		return err
	}

	if action.Body != nil {
		err = writeClaimsToBody(r, claims, &action)
		if err != nil {
			return err
		}
	}
	if action.Query != nil {
		err = writeClaimsToQuery(r, claims, &action)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns the jwt claim object if valid
func getJWTClaims(r *http.Request, action *actions.JWTExtract) (jwt.MapClaims, error) {
	token, err := getAuthToken(r)
	if err != nil {
		return nil, err
	}

	s, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(action.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := s.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Cannot map jwt claims")
	}

	return claims, nil
}

// Fetches the auth token from the authorization header
// disacrds the bearer prefix
func getAuthToken(r *http.Request) (string, error) {
	requestToken := r.Header.Get("Authorization")
	splitToken := strings.Split(requestToken, "Bearer ")
	if len(splitToken) < 2 {
		return "", errors.New("missing auth token")
	}
	return splitToken[1], nil
}

// Writes the action's parameters from the jwt claims to the request body
// Overrides existing keys with duplicate names
func writeClaimsToBody(r *http.Request, claims jwt.MapClaims, action *actions.JWTExtract) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	bodyFormatted := make(map[string]interface{})
	err = json.Unmarshal(b, &bodyFormatted)
	if err != nil {
		return err
	}

	for bodyKey, claimsKey := range action.Body {
		if claimsValue, ok := claims[claimsKey]; ok {
			bodyFormatted[bodyKey] = claimsValue
		}
	}

	b, err = json.Marshal(bodyFormatted)
	if err != nil {
		return err
	}

	r.ContentLength = int64(len(b))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return nil
}

// Writes the actions' parameters from the jwt claims to the request query
// Overrides existing keys with duplicate names
func writeClaimsToQuery(r *http.Request, claims jwt.MapClaims, action *actions.JWTExtract) error {
	q := r.URL.Query()
	for queryKey, claimsKey := range action.Query {
		if claimValue, ok := claims[claimsKey]; ok {
			if claimValueStr, ok := claimValue.(string); ok {
				q.Add(queryKey, claimValueStr)
			}
		}
	}

	r.URL.RawQuery = q.Encode()

	return nil
}
