package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/codegangsta/cli"
	"github.com/wsxiaoys/terminal/color"
)

var (
	ErrItemNotFound = errors.New("specfied item not found")
	ErrInvalidIput  = errors.New("invalid input")
	ErrInvalidTags  = errors.New("invalid tags format")
)

var getCommand = cli.Command{
	Name:   "get",
	Usage:  "retrieve items from pocket",
	Action: getAction,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "count,c",
			Usage: "only return count number of numbers",
		},
		cli.StringFlag{
			Name:  "tag,t",
			Value: "",
			Usage: "return items tagged with specified tag,or \"_untagged\" for untagged items",
		},
		cli.StringFlag{
			Name:  "state,s",
			Value: "unread",
			Usage: "return unread / archive / all items",
		},
		cli.BoolFlag{
			Name:  "favorite,f",
			Usage: "only return favorited items",
		},
		cli.StringFlag{
			Name:  "mime,m",
			Value: "",
			Usage: "only return specified item type : article/video/image",
		},
		cli.StringFlag{
			Name:  "order,o",
			Value: "newest",
			Usage: "return item in specified order : newest/oldest/title/site",
		},
		cli.StringFlag{
			Name:  "domain,d",
			Value: "",
			Usage: "only return items from a particular domain",
		},
	},
}

func getAction(c *cli.Context) {
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	count := c.Int("count")

	tag := c.String("tag")
	state := c.String("state")
	fav := c.Bool("favorite")
	contentType := c.String("mime")
	order := c.String("order")
	domain := c.String("domain")

	values := map[string]string{}
	values["consumer_key"] = cfg.ConsumerKey
	values["access_token"] = cfg.AccessToken
	values["detailType"] = "simple"

	if count > 0 {
		values["count"] = strconv.Itoa(count)
	}
	if tag != "" {
		values["tag"] = tag
	}

	if domain != "" {
		values["domain"] = domain
	}

	if contentType != "" {
		if contains([]string{"image", "article", "video"}, contentType) {
			values["contentType"] = contentType
		} else {
			logger.Fatal("content type args error,must be one of 'article/video/image'")
		}
	}

	if contains([]string{"unread", "archive", "all"}, state) {
		values["state"] = state
	} else {
		logger.Fatal("item state args error ,must be one of 'unread/archive/all'")
	}

	if contains([]string{"newest", "oldest", "title", "site"}, order) {
		values["sort"] = order
	} else {
		logger.Fatal("items order args error,must be one of 'newest/oldest/title/site'")
	}

	if fav == true {
		values["favorite"] = "1"
	}

	b, err := json.Marshal(values)
	if err != nil {
		logger.Fatal(err)
	}

	resp, err := doRequest("/v3/get", "POST", 200, b)
	if err != nil {
		logger.Fatal(err)
	}

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}

	store, err := formatItems(json)
	if err != nil {
		logger.Fatal(err)
	}

	printActions(store)

}

func printActions(store *map[string]Meta) {
	reader := bufio.NewReader(os.Stdin)

	for {
		color.Println("@r\nActions:")
		color.Println("@b[v] view in default browser (eg: v 1)")
		color.Println("@b[c] view in terminal(need w3m)")
		color.Println("@b[a] archive")
		color.Println("@b[f] favorite")
		color.Println("@b[d] delete")
		color.Println("@b[t] tags.set new tags for this item")
		color.Println("@b[q] quit")
		color.Print("@b====> ")
		action, _ := reader.ReadString('\n')
		if strings.Split(strings.TrimSpace(action), " ")[0] == "q" {
			break
		} else {
			err := processActions(strings.TrimSpace(action), store)
			if err != nil {
				fmt.Println(err)
			}

		}
	}

}

func processActions(action string, store *map[string]Meta) error {

	s := strings.Split(action, " ")
	if len(s) < 2 {
		return ErrInvalidIput
	}

	meta, ok := (*store)[s[1]]
	if !ok {
		return ErrItemNotFound
	}

	switch s[0] {
	case "v":
		if err := viewInBrowser(meta); err != nil {
			return err
		}

	case "a":
		if err := archive(meta); err != nil {
			return err
		} else {
			fmt.Println("Item archived.")
		}
	case "c":
		if err := cat(meta); err != nil {
			return err
		}
	case "f":
		if err := favorite(meta); err != nil {
			return err
		} else {
			fmt.Println("Item mark as favorite.")
		}
	case "d":
		if err := delete(meta); err != nil {
			return err
		} else {
			fmt.Println("Item deleted")
		}
	case "t":
		if len(s) != 3 {
			return ErrInvalidTags
		}
		if err := tags(meta, s[2]); err != nil {
			return err
		} else {
			fmt.Println("Item tags updated.")
		}

	}
	return nil
}

func formatItems(js *simplejson.Json) (*map[string]Meta, error) {
	result, err := js.Map()
	if err != nil {
		return nil, err
	}
	listValue := result["list"]
	items := reflect.ValueOf(listValue).Interface().(map[string]interface{})

	count := 0
	store := make(map[string]Meta)
	for key, value := range items {
		count = count + 1
		color.Printf("@m[%v]", count)
		color.Printf("@y[%s]:", key)
		//fmt.Printf("[%s]:", key)
		vs := reflect.ValueOf(value).Interface().(map[string]interface{})
		title := reflect.ValueOf(vs["given_title"]).Interface().(string)
		url := reflect.ValueOf(vs["given_url"]).Interface().(string)
		color.Printf("@c[%s] ", title)
		color.Printf("@g%v\n", url)

		store[strconv.Itoa(count)] = Meta{Url: url, Title: title, Id: key}
	}
	return &store, nil
}
