package twitter_internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (r *Twitter) GetUserByName(name string) (*User, error) {
	r.refreshGuestToken()

	val := fmt.Sprintf(`{"screen_name":%q,"withSafetyModeUserFields":true,"withSuperFollowsUserFields":true}`, name)
	uri := fmt.Sprintf("https://twitter.com/i/api/graphql/B-dCk4ph5BZ0UReWK590tw/UserByScreenName?variables=%s", url.QueryEscape(val))
	resp := new(getUserByNameResp)
	req := r.req(http.MethodGet, uri)
	text, err := req.Text()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(text), resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal user %s fail", text)
	}
	return resp.Data.User.Result, nil
}

type getUserByNameResp struct {
	Data struct {
		User struct {
			Result *User `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type User struct {
	ID     string `json:"id"`
	RestID string `json:"rest_id"`
	Legacy struct {
		CreatedAt            string   `json:"created_at"`
		Description          string   `json:"description"`
		Name                 string   `json:"name"`
		PinnedTweetIdsStr    []string `json:"pinned_tweet_ids_str"`
		ProfileBannerURL     string   `json:"profile_banner_url"`
		ProfileImageURLHTTPS string   `json:"profile_image_url_https"`
		ScreenName           string   `json:"screen_name"`
		URL                  string   `json:"url"`
	} `json:"legacy"`
}
