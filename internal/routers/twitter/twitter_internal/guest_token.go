package twitter_internal

import (
	"log"
	"net/http"
)

func (r *Twitter) refreshGuestToken() {
	r.once.Do(func() {
		resp := new(refreshGuestTokenResp)
		err := r.req(http.MethodPost, "https://api.twitter.com/1.1/guest/activate.json").Unmarshal(resp)
		if err != nil {
			panic(err)
		}
		r.guestToken = resp.GuestToken
		log.Printf("[twitter] refresh gust-token: %s", r.guestToken)
	})
}

type refreshGuestTokenResp struct {
	GuestToken string `json:"guest_token"`
}
