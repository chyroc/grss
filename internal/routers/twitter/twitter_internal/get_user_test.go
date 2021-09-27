package twitter_internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetUser(t *testing.T) {
	as := assert.New(t)

	user, err := New().GetUserByName("awscloud")
	as.Nil(err)
	as.NotNil(user)
	as.Equal("Amazon Web Services", user.Legacy.Name)
	as.Equal("66780587", user.RestID)
}
