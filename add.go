package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
)

type Item struct {
	ItemId         string `json:"item_id,omitempty"`
	NormalUrl      string `json:"normal_url,omitempty"`
	ResolvedId     string `json:"resolved_id,omitempty"`
	ResolvedUrl    string `json:"resolved_url,omitempty"`
	DomainId       string `json:"domain_id,omitempty"`
	OriginDomainId string `json:"origin_domain_id,omitempty"`
	ResponseCode   string `json:"response_code,omitempty"`
	MimeType       string `json:"mime_type,omitempty"`
	ContentLength  string `json:"content_length,omitempty"`
	Encoding       string `json:"encoding,omitempty"`
	DateResolved   string `json:"date_resolved,omitempty"`
	DatePublished  string `json:"date_published,omitempty"`
	Title          string `json:"title,omitempty"`
	Excerpt        string `json:"excerpt,omitempty"`
	WordCount      string `json:"word_count,omitempty"`
	HasImage       string `json:"has_image,omitempty"`
	HasVideo       string `json:"has_video,omitempty"`
	IsIndex        string `json:"is_index,omitempty"`
	IsArticle      string `json:"is_article,omitempty"`
}

type AddResp struct {
	Item   Item `json:"item,omitempty"`
	Status int  `json:"status,omitempty"`
}

var addCommand = cli.Command{
	Name:   "add",
	Usage:  "add an item to pocket",
	Action: addAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "url,u",
			Value: "",
			Usage: "The url you want to add to pocket",
		},
		cli.StringFlag{
			Name:  "name,n",
			Value: "",
			Usage: "the title of this article if it don't have",
		},
		cli.StringFlag{
			Name:  "tags,t",
			Value: "",
			Usage: "A comma-separated list of tags to apply to the item",
		},
	},
}

func addAction(c *cli.Context) {
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	url := c.String("url")
	title := c.String("name")
	tags := c.String("tags")

	if url == "" {
		logger.Fatal("you must specify the url of this item")
	}
	add, err := add(url, title, tags, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	printAddInfo(&add.Item)
}

func add(url, title, tags string, cfg *PocketConfig) (*AddResp, error) {
	values := map[string]string{}
	values["consumer_key"] = cfg.ConsumerKey
	values["access_token"] = cfg.AccessToken
	values["url"] = url
	if title != "" {
		values["title"] = title
	}
	if tags != "" {
		values["tags"] = tags
	}
	b, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	resp, err := doRequest("/v3/add", "POST", 200, b)
	if err != nil {
		return nil, err
	}

	var add *AddResp
	if err := json.NewDecoder(resp.Body).Decode(&add); err != nil {
		return nil, err
	}
	return add, nil

}

func printAddInfo(item *Item) error {
	fmt.Printf("ResolvedUrl : %v\n", item.ResolvedUrl)
	fmt.Printf("Title       : %v\n", item.Title)
	fmt.Printf("MimeType    : %v\n", item.MimeType)
	fmt.Printf("Encoding    : %v\n", item.Encoding)
	return nil
}
