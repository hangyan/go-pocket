package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func buildUrl(path string) string {
	return fmt.Sprintf("%s%s", POCKET_URL, path)
}

func doRequest(path string, method string, expectedStatus int, b []byte) (*http.Response, error) {
	url := buildUrl(path)
	buf := bytes.NewBuffer(b)

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/json;charset=UTF8")
	req.Header.Add("X-Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedStatus {
		c, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return resp, errors.New(string(c) + "==>" + resp.Header.Get("X-Error"))
	}
	return resp, nil
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}
