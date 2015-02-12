package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"reflect"

	"github.com/bitly/go-simplejson"
	"github.com/skratchdot/open-golang/open"
)

var (
	ErrItemArchive  = errors.New("item archive failed")
	ErrItemFavorite = errors.New("item favorite failed")
	ErrItemDelete   = errors.New("item delete failed")
)

func viewInBrowser(meta Meta) error {
	err := open.Run(meta.Url)
	if err != nil {
		return err
	}

	return nil
}

func buildGetUrl(actions string) (string, error) {

	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/v3/send?actions=%s&access_token=%s&consumer_key=%s", actions, cfg.AccessToken, cfg.ConsumerKey), nil

}

func favorite(meta Meta) error {

	favoriteAction := fmt.Sprintf("[{\"action\":\"favorite\",\"item_id\":\"%s\"}]", meta.Id)
	request, err := buildGetUrl(favoriteAction)
	if err != nil {
		return err
	}
	resp, err := http.Get(buildUrl(request))
	if err != nil {
		return err
	}
	result, err := checkActionResult(resp.Body)
	if err != nil {
		return err
	}
	if result == false {
		return ErrItemFavorite
	}
	return nil
}

func checkActionResult(body io.ReadCloser) (bool, error) {
	js, err := simplejson.NewFromReader(body)
	if err != nil {
		return false, err
	}

	result, err := js.Map()
	if err != nil {
		return false, err
	}

	status, err := reflect.ValueOf(result["status"]).Interface().(json.Number).Int64()
	if err != nil {
		return false, err
	}
	if status == 1 {
		return true, nil
	} else {
		return false, nil
	}

}

func delete(meta Meta) error {
	deleteAction := fmt.Sprintf("[{\"action\":\"delete\",\"item_id\":\"%s\"}]", meta.Id)
	request, err := buildGetUrl(deleteAction)
	if err != nil {
		return err
	}
	resp, err := http.Get(buildUrl(request))
	if err != nil {
		return err
	}

	result, err := checkActionResult(resp.Body)
	if err != nil {
		return err
	}
	if !result {
		return ErrItemDelete
	}
	return nil
}

func archive(meta Meta) error {
	archiveAction := fmt.Sprintf("[{\"action\":\"archive\",\"item_id\":\"%s\"}]", meta.Id)
	request, err := buildGetUrl(archiveAction)
	if err != nil {
		return err
	}
	resp, err := http.Get(buildUrl(request))
	if err != nil {
		return err
	}

	result, err := checkActionResult(resp.Body)
	if err != nil {
		return err
	}
	if !result {
		return ErrItemArchive
	}

	return nil
}

func cat(meta Meta) error {
	_, err := exec.LookPath("w3m")
	if err != nil {
		return err
	} else {
		out, err := exec.Command("w3m", "-dump", meta.Url).Output()
		if err != nil {
			return err
		} else {
			fmt.Printf("-----------------------------------\n%s", out)
		}
	}

	return nil

}
