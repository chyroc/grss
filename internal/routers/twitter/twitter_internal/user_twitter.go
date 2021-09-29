package twitter_internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/helper"
)

func (r *Twitter) GetUserTwitter(userRestID string) ([]*GetUserTwitterRespEntry, error) {
	r.refreshGuestToken()

	uri := "https://twitter.com/i/api/graphql/Lya9A5YxHQxhCQJ5IPtm7A/UserTweets"
	val := fmt.Sprintf(`{"userId":%q,"count":100,"withTweetQuoteCount":true,"includePromotedContent":true,"withSuperFollowsUserFields":true,"withUserResults":true,"withBirdwatchPivots":false,"withReactionsMetadata":false,"withReactionsPerspective":false,"withSuperFollowsTweetFields":true,"withVoice":true}`, userRestID)
	req := r.req(http.MethodGet, uri).WithQuery("variables", val)

	text, err := req.Text()
	if err != nil {
		return nil, err
	}
	resp := new(getUserTwitterResp)
	err = json.Unmarshal([]byte(text), resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal user %s fail", text)
	}
	if len(resp.Data.User.Result.Timeline.Timeline.Instructions) == 0 {
		return nil, nil
	}

	// spew.Dump(text)

	resp2 := []*GetUserTwitterRespEntry{}
	err = lambda.New(resp.Data.User.Result.Timeline.Timeline.Instructions[0].Entries).FilterArray(func(idx int, obj interface{}) bool {
		return obj.(*GetUserTwitterRespEntry).Content.EntryType != "TimelineTimelineCursor"
	}).ToList(&resp2)
	return resp2, err
}

type getUserTwitterResp struct {
	Data struct {
		User struct {
			Result struct {
				Typename string `json:"__typename"`
				Timeline struct {
					Timeline struct {
						Instructions []struct {
							Entries []*GetUserTwitterRespEntry `json:"entries,omitempty"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type getUserTwitterRespLegacy struct {
	CreatedAt           string `json:"created_at"`
	DefaultProfile      bool   `json:"default_profile"`
	DefaultProfileImage bool   `json:"default_profile_image"`
	Description         string `json:"description"`
	Entities            struct {
		Description struct {
			Urls []interface{} `json:"urls"`
		} `json:"description"`
		URL struct {
			Urls []struct {
				DisplayURL  string `json:"display_url"`
				ExpandedURL string `json:"expanded_url"`
				URL         string `json:"url"`
				Indices     []int  `json:"indices"`
			} `json:"urls"`
		} `json:"url"`
	} `json:"entities"`
	FastFollowersCount      int           `json:"fast_followers_count"`
	FavouritesCount         int           `json:"favourites_count"`
	FollowersCount          int           `json:"followers_count"`
	FriendsCount            int           `json:"friends_count"`
	HasCustomTimelines      bool          `json:"has_custom_timelines"`
	IsTranslator            bool          `json:"is_translator"`
	ListedCount             int           `json:"listed_count"`
	Location                string        `json:"location"`
	MediaCount              int           `json:"media_count"`
	Name                    string        `json:"name"`
	NormalFollowersCount    int           `json:"normal_followers_count"`
	PinnedTweetIdsStr       []string      `json:"pinned_tweet_ids_str"`
	ProfileBannerURL        string        `json:"profile_banner_url"`
	ProfileImageURLHTTPS    string        `json:"profile_image_url_https"`
	ProfileInterstitialType string        `json:"profile_interstitial_type"`
	Protected               bool          `json:"protected"`
	ScreenName              string        `json:"screen_name"`
	StatusesCount           int           `json:"statuses_count"`
	TranslatorType          string        `json:"translator_type"`
	URL                     string        `json:"url"`
	Verified                bool          `json:"verified"`
	WithheldInCountries     []interface{} `json:"withheld_in_countries"`
}

type GetUserTwitterRespEntry struct {
	EntryID   string `json:"entryId"`
	SortIndex string `json:"sortIndex"`
	Content   struct {
		EntryType   string `json:"entryType"`
		ItemContent struct {
			ItemType         string                          `json:"itemType"`
			TweetResults     getUserTwitterRespTwitterResult `json:"tweet_results"`
			TweetDisplayType string                          `json:"tweetDisplayType"`
		} `json:"itemContent"`
		Items []struct {
			EntryID string `json:"entryId"`
			Item    struct {
				ItemContent struct {
					ItemType string `json:"itemType"`
					Topic    struct {
						Description   string `json:"description"`
						Following     bool   `json:"following"`
						IconURL       string `json:"icon_url"`
						ID            string `json:"id"`
						Name          string `json:"name"`
						NotInterested bool   `json:"not_interested"`
					} `json:"topic"`
					TopicFunctionalityType string `json:"topicFunctionalityType"`
					TopicDisplayType       string `json:"topicDisplayType"`
				} `json:"itemContent"`
				FeedbackInfo struct {
					FeedbackKeys     []string `json:"feedbackKeys"`
					FeedbackMetadata string   `json:"feedbackMetadata"`
				} `json:"feedbackInfo"`
				ClientEventInfo struct {
					Component string `json:"component"`
					Element   string `json:"element"`
					Details   struct {
						TimelinesDetails struct {
							ControllerData string `json:"controllerData"`
						} `json:"timelinesDetails"`
					} `json:"details"`
				} `json:"clientEventInfo"`
			} `json:"item"`
		} `json:"items"`
		Metadata struct {
			GridCarouselMetadata struct {
				NumRows int `json:"numRows"`
			} `json:"gridCarouselMetadata"`
		} `json:"metadata"`
		DisplayType string `json:"displayType"`
		Header      struct {
			DisplayType   string `json:"displayType"`
			Text          string `json:"text"`
			SocialContext struct {
				Type        string `json:"type"`
				ContextType string `json:"contextType"`
				Text        string `json:"text"`
			} `json:"socialContext"`
			Sticky bool `json:"sticky"`
		} `json:"header"`
		ClientEventInfo struct {
			Component string `json:"component"`
		} `json:"clientEventInfo"`
		Value               string `json:"value"`
		CursorType          string `json:"cursorType"`
		StopOnEmptyResponse bool   `json:"stopOnEmptyResponse"`
	} `json:"content,omitempty"`
}

type getUserTwitterRespTwitterResult struct {
	Result struct {
		Typename string `json:"__typename"`
		RestID   string `json:"rest_id"`
		Core     struct {
			UserResults struct {
				Result struct {
					Typename                   string                   `json:"__typename"`
					ID                         string                   `json:"id"`
					RestID                     string                   `json:"rest_id"`
					AffiliatesHighlightedLabel struct{}                 `json:"affiliates_highlighted_label"`
					Legacy                     getUserTwitterRespLegacy `json:"legacy"`
					SuperFollowEligible        bool                     `json:"super_follow_eligible"`
					SuperFollowedBy            bool                     `json:"super_followed_by"`
					SuperFollowing             bool                     `json:"super_following"`
				} `json:"result"`
			} `json:"user_results"`
		} `json:"core"`
		Legacy struct {
			CreatedAt         string `json:"created_at"`
			ConversationIDStr string `json:"conversation_id_str"`
			DisplayTextRange  []int  `json:"display_text_range"`
			Entities          struct {
				UserMentions []struct {
					IDStr      string `json:"id_str"`
					Name       string `json:"name"`
					ScreenName string `json:"screen_name"`
					Indices    []int  `json:"indices"`
				} `json:"user_mentions"`
				Urls     []interface{} `json:"urls"`
				Hashtags []interface{} `json:"hashtags"`
				Symbols  []interface{} `json:"symbols"`
				Media    []struct {
					MediaURLHTTPS string `json:"media_url_https"`
					Type          string `json:"type"`
					URL           string `json:"url"`
				} `json:"media"`
			} `json:"entities"`
			FavoriteCount         int                              `json:"favorite_count"`
			Favorited             bool                             `json:"favorited"`
			FullText              string                           `json:"full_text"`
			IsQuoteStatus         bool                             `json:"is_quote_status"`
			Lang                  string                           `json:"lang"`
			QuoteCount            int                              `json:"quote_count"`
			ReplyCount            int                              `json:"reply_count"`
			RetweetCount          int                              `json:"retweet_count"`
			Retweeted             bool                             `json:"retweeted"`
			Source                string                           `json:"source"`
			UserIDStr             string                           `json:"user_id_str"`
			IDStr                 string                           `json:"id_str"`
			RetweetedStatusResult *getUserTwitterRespTwitterResult `json:"retweeted_status_result"`
		} `json:"legacy"`
	} `json:"result"`
}

func (r *GetUserTwitterRespEntry) RetweetedResult() *getUserTwitterRespTwitterResult {
	return r.Content.ItemContent.TweetResults.Result.Legacy.RetweetedStatusResult
}

func (r *GetUserTwitterRespEntry) IsRetwitter() bool {
	return r.RetweetedResult() != nil
}

func (r *GetUserTwitterRespEntry) OriginHtml() string {
	text := helper.ToHtml(r.Content.ItemContent.TweetResults.Result.Legacy.FullText)
	for _, v := range r.Content.ItemContent.TweetResults.Result.Legacy.Entities.Media {
		if v.Type != "photo" {
			continue
		}
		if strings.Contains(text, v.URL) {
			text=strings.ReplaceAll(text, v.URL, fmt.Sprintf("<div><img src=%q alt=%q/></div>", v.MediaURLHTTPS, v.URL))
		}
	}
	return text
}
