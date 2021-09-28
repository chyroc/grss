package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToTitleString(t *testing.T) {
	assert.Equal(t, "$酒鬼酒(SZ000799)$", ToTitleText(`<a href="http://xueqiu.com/S/SZ000799" target="_blank">$酒鬼酒(SZ000799)$</a>`, 100, " ..."))
	assert.Equal(t, "视频B站同步发了，欢迎关注。是的，我也有B站账号的  只要5分钟，在云服务器上一键安装180个优质云应用", ToTitleText(`视频B站同步发了，欢迎关注。是的，我也有B站账号的 <span class="url-icon"><img alt=[允悲] src="xxx" /></span> <a data-url="http://t.cn/A6MtWp6T" href="http://t.cn/A6MtWp6T" data-hide=""><span class='url-icon'><img style='width: 1rem;height: 1rem' src='https://h5.sinaimg.cn/upload/2015/09/25/3/timeline_card_small_web_default.png'></span><span class="surl-text">只要5分钟，在云服务器上一键安装180个优质云应用</span></a>`, 100, " ..."))
}
