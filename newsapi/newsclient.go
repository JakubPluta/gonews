package newsapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

var NewsApiBaseURL string = "https://newsapi.org/v2/"

func (a *Article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func (c *Client) FetchEverything(q string, page string) (*Results, error) {
	endpoint := fmt.Sprintf("%severything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", NewsApiBaseURL, url.QueryEscape(q), c.PageSize, page, c.key)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code %d from %s with body %s", resp.StatusCode, endpoint, string(body))
	}
	res := &Results{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func NewClient(httpClient *http.Client, apiKey string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}
	return &Client{
		http:     httpClient,
		key:      apiKey,
		PageSize: pageSize,
	}
}
