package authcli

import (
	"bytes"
	"encoding/json"
	"errors"

	"net/http"
)

// RefreshToken sends a refresh token request to IAM and returns a new token
func (cli *Cli) RefreshToken(refreshToken string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{"refreshToken": refreshToken})
	if err != nil {
		return "", err
	}

	resp, err := cli.httpClient.Post(cli.AuthNSvrURL+"api/v1/auth/refresh", "application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to refresh token")
	}

	var res struct {
		Token string `json:"token"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Token, nil
}
