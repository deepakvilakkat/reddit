// Package reddit implements a basic client for Reddit API.
package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Item describes a Reddit item.
type Item struct {
	Title    string
	URL      string
	Comments int `json:"num_comments"`
}

type response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

// Get fetches the most recent Items posted to the specified subreddit.
func Get(reddit string) ([]Item, error) {
	url := fmt.Sprintf("http://reddit.com/r/%s.json", reddit)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	r := new(response)
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return nil, err
	}

	items := make([]Item, len(r.Data.Children))
	for i, child := range r.Data.Children {
		items[i] = child.Data
	}
	return items, nil
}

func (i Item) String() string {
	com := ""
	switch i.Comments {
	case 0:
		com = "No Comments"
	case 1:
		com = "1 Comment"
	default:
		com = fmt.Sprintf("%d Comments", i.Comments)
	}
	return fmt.Sprintf("%s - %s ( %s )", i.Title, i.URL, com)
}
