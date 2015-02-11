package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/codegangsta/cli"
	"strconv"
)

var getCommand = cli.Command{
	Name:   "get",
	Usage:  "retrieve items from pocket",
	Action: getAction,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "count,c",
			Value: 1,
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
			Usage: "only return specified item type : article/video/image",
		},
		cli.StringFlag{
			Name:  "order,o",
			Value: "newest",
			Usage: "return item in specified order : newest/oldest/title/site",
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

	values := map[string]string{}
	values["consumer_key"] = cfg.ConsumerKey
	values["access_token"] = cfg.AccessToken

	if count > 0 {
		values["count"] = strconv.Itoa(count)
	}
	if tag != "" {
		values["tag"] = tag
	}

	values["detailType"] = "simple"

	if contentType != "" {
		if contains([]string{"image", "article", "video"}, contentType) {
			values["contentType"] = contentType
		} else {
			logger.Fatal("content type args error,must be one of 'article/video/image'")
		}
	}

	if state == "unread" || state == "archive" || state == "all" {
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
	} else {
		values["favorite"] = "0"
	}

	b, err := json.Marshal(values)
	if err != nil {
		logger.Fatal(err)
	}

	resp, err := doRequest("/v3/get", "POST", 200, b)
	if err != nil {
		logger.Fatal(err)
	}

	// var items *GetResp
	// if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
	//	logger.Fatal(err)
	//}
	// printGetItems(items)

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(json)

}
