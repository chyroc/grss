package twitter_user_origin

import (
	"fmt"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
	"github.com/chyroc/grss/internal/routers/twitter/twitter_internal"
)

// args: map[string]string{"uid": "woshisuchang"}
func New(args map[string]string) (*fetch.Source, error) {
	uid := args["uid"]
	user, err := twitter_internal.Ins.GetUserByName(uid)
	if err != nil {
		return nil, err
	}
	return &fetch.Source{
		Title: fmt.Sprintf("Twitter - %s Origin Twitter", user.Legacy.Name),
		Link:  fmt.Sprintf("https://twitter.com/%s/", uid),

		Fetch: func() (interface{}, error) {
			entry, err := twitter_internal.Ins.GetUserTwitter(user.RestID)
			return entry, err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj.([]*twitter_internal.GetUserTwitterRespEntry)).FilterArray(func(idx int, obj interface{}) bool {
				return !obj.(*twitter_internal.GetUserTwitterRespEntry).IsRetwitter()
			}).MapArrayAsync(func(idx int, obj interface{}) interface{} {
				entry := obj.(*twitter_internal.GetUserTwitterRespEntry)
				link := fmt.Sprintf("https://twitter.com/%s/status/%s", uid, entry.EntryID)
				pubTime, _ := time.Parse("Mon Jan 02 15:04:05 -0700 2006", entry.Content.ItemContent.TweetResults.Result.Legacy.CreatedAt)

				return &fetch.Item{
					Title:       helper.InterceptString(entry.Content.ItemContent.TweetResults.Result.Legacy.FullText, 100, " ..."),
					Link:        link,
					Description: entry.Content.ItemContent.TweetResults.Result.Legacy.FullText,
					Author:      user.Legacy.Name,
					PubDate:     pubTime,
				}
			}).ToList(&resp)
			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}
