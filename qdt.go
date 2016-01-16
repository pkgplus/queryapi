package queryapi

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	STR_REG_QDT = `(?s)class="place".*?<strong>(截至.*?)<\/strong>`
)

var RET_QDT *regexp.Regexp = regexp.MustCompile(STR_REG_QDT)

func QueryQDT(cardno string) (string, error) {
	url := fmt.Sprintf("http://www.qdtcn.com/qdtnet/query.do?card_no=%s", cardno)
	res, err := http.Post(url, "text/html;charset=utf-8", bytes.NewBuffer([]byte("")))
	if err != nil {
		return "查询失败(POST_ERR)", err
	}
	defer res.Body.Close()

	body_bytes, err_read := ioutil.ReadAll(res.Body)
	if err_read != nil {
		return "查询失败(READ_BODY_ERR)", err_read
	}
	matched_strs := RET_QDT.FindStringSubmatch(string(body_bytes))
	if len(matched_strs) >= 1 {
		return matched_strs[1], nil
	} else {
		return "截取失败(READ_BODY_ERR)", errors.New("正则截取余额失败！")
	}
}
