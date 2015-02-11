package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/skratchdot/open-golang/open"
	"os"
)

type (
	AuthResp struct {
		Code  string `json:"code,omitempty"`
		State string `json:"state,omitempty"`
	}

	AccessResp struct {
		AccessToken string `json:"access_token,omitempty"`
		Username    string `json:"username,omitempty"`
	}
)

var authCommand = cli.Command{
	Name:   "auth",
	Usage:  "Obtain a request token",
	Action: authAction,
}

func authAction(c *cli.Context) {
	resp, err := auth()
	if err != nil {
		logger.Fatal(err)
	}

	err = open.Run(POCKET_URL + "/auth/authorize?request_token=" + resp.Code + "&redirect_uri=" + REDIRECT_URI)
	if err != nil {
		logger.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Type anything and this will be done => ")
	_, _ = reader.ReadString('\n')
	fmt.Println("That's right !You are awesome!")
	access, err := access(resp.Code)
	if err != nil {
		logger.Fatal(err)
	}

	cfg := &PocketConfig{
		ConsumerKey: CONSUMER_KEY,
		Username:    access.Username,
		AccessToken: access.AccessToken,
	}

	if err := saveConfig(cfg); err != nil {
		logger.Fatal(err)
	}

}

func auth() (*AuthResp, error) {
	values := map[string]string{}
	values["consumer_key"] = CONSUMER_KEY
	values["redirect_uri"] = REDIRECT_URI
	b, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest("/v3/oauth/request", "POST", 200, b)
	if err != nil {
		return nil, err
	}
	var arp *AuthResp
	if err := json.NewDecoder(resp.Body).Decode(&arp); err != nil {
		return nil, err
	}

	return arp, nil

}

func access(code string) (*AccessResp, error) {
	values := map[string]string{}
	values["consumer_key"] = CONSUMER_KEY
	values["code"] = code
	b, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest("/v3/oauth/authorize", "POST", 200, b)
	if err != nil {
		return nil, err
	}

	var access *AccessResp
	if err := json.NewDecoder(resp.Body).Decode(&access); err != nil {
		return nil, err
	}
	return access, nil
}
