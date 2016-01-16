package queryapi

import (
	"encoding/json"
	"errors"
	"fmt"
	json_sample "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

type TuringNews struct {
	Code int    `json:"code"`
	Text string `json:"text"`
	List []struct {
		Article   string `json:"article"`
		Source    string `json:"source"`
		Icon      string `json:"icon"`
		Detailurl string `json:"detailurl"`
	} `json:"list"`
}

type TuringCookBook struct {
	Code int    `json:"code"`
	Text string `json:"text"`
	List []struct {
		Name      string `json:"name"`
		Icon      string `json:"icon"`
		Info      string `json:"info"`
		Detailurl string `json:"detailurl"`
	} `json:"list"`
}

const (
	FMT_RUL_TURING = `http://www.tuling123.com/openapi/api?userid=%s&key=%s&info=%s`
)

func QueryTURING(userid, text, api_key string) (interface{}, error) {
	if api_key == "" {
		return "", errors.New("Not Found TURING APIKEY!")
	}

	get_url := fmt.Sprintf(FMT_RUL_TURING, userid, api_key, text)
	resp, err2 := http.Get(get_url)
	if err2 != nil {
		return "", err2
	}
	defer resp.Body.Close()

	resp_body, err_tmp := ioutil.ReadAll(resp.Body)
	if err_tmp != nil {
		return "", err_tmp
	}

	ret_json, j_err := json_sample.NewJson(resp_body)
	if j_err != nil {
		return "", j_err
	}

	json_map, map_err := ret_json.Map()
	if map_err != nil {
		return "", map_err
	}

	code, found := json_map["code"]
	if !found {
		return "", errors.New("Not found code!")
	}
	code_i, int_err := strconv.Atoi(reflect.ValueOf(code).String())
	if int_err != nil {
		return "", errors.New("Parse code to int failed!" + int_err.Error())
	}

	switch code_i {
	case 100000:
		return ParseText(json_map)
	case 200000:
		return ParseImage(json_map)
	case 302000:
		return ParseNews(resp_body)
	case 308000:
		return ParseCB(resp_body)
	}

	return "", nil
}

func ParseCB(body_bytes []byte) ([]News, error) {
	q_news := make([]News, 0)
	t_news := &TuringNews{}
	err := json.Unmarshal(body_bytes, t_news)
	if err != nil {
		return q_news, err
	}

	for _, news := range t_news.List {
		q_news = append(q_news, News{news.Article,
			news.Source,
			news.Icon,
			news.Detailurl})
	}

	return q_news, nil
}

func ParseNews(body_bytes []byte) ([]News, error) {
	q_news := make([]News, 0)
	t_news := &TuringNews{}
	err := json.Unmarshal(body_bytes, t_news)
	if err != nil {
		return q_news, err
	}

	for _, news := range t_news.List {
		q_news = append(q_news, News{news.Article,
			news.Source,
			news.Icon,
			news.Detailurl})
	}

	return q_news, nil
}

func ParseImage(json_map map[string]interface{}) (string, error) {
	text_v, found := json_map["text"]
	if !found {
		return "", errors.New("Not found text!")
	}

	text, ok := text_v.(string)
	if !ok {
		return "", errors.New("Parse text to string error!")
	}

	url_v, found := json_map["url"]
	if !found {
		return "", errors.New("Not url text!")
	}

	url_str, ok := url_v.(string)
	if !ok {
		return "", errors.New("Parse url to string error!")
	}

	return fmt.Sprintf("%s\n<a href=\"%s\">点击查看</a>", text, url_str), nil
}

func ParseText(json_map map[string]interface{}) (string, error) {
	text_v, found := json_map["text"]
	if !found {
		return "", errors.New("Not found text!")
	}

	text, ok := text_v.(string)
	if !ok {
		return "", errors.New("Parse text to string error!")
	}

	return text, nil
}
