package queryapi

import (
	"regexp"
)

const (
	WEIZHANG_URL = `http://www.cheshouye.com/api/weizhang/`
	HM_REG_STR   = `(?s)hm.baidu.com\/h.js%3F(\w+)`
)

var REG_HM *regexp.Regexp

func init() {
	REG_HM = regexp.MustCompile(HM_REG_STR)
}
