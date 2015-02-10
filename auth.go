package main

import (
	"encoding/json"
	"github.com/codegangsta/cli"
)

type (
	AuthResp struct {
		code  string `json:"code"`
		state string `json:"state"`
	}
)

var authCommand = cli.Command{
	Name:   "auth",
	Usage:  "Obtain a request token",
	Action: authAction,
}

func authAction(c *cli.Context) {
	resp, err := Auth()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("hehe : %s", resp.code)

}

func Auth() (*AuthResp, error) {
	values := map[string]string{}
	values["consumer_key"] = CONSUMER_KEY
	values["redirect_uri"] = REDIRECT_URI
	b, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest("/oauth/request", "POST", 200, b)
	if err != nil {
		return nil, err
	}

	var ap *AuthResp
	if err := json.NewDecoder(resp.Body).Decode(&ap); err != nil {
		return nil, err
	}

	return ap, nil

}
