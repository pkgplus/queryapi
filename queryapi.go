package queryapi

import (
	"errors"
	"strings"
)

type API struct {
	KEY map[string]string
}

type News struct {
	Article   string `json:"article"`
	Source    string `json:"source"`
	Icon      string `json:"icon"`
	Detailurl string `json:"detailurl"`
}

func (a *API) SetKey(k_type, k_detail string) {
	if len(a.KEY) == 0 {
		a.KEY = make(map[string]string)
	}
	a.KEY[k_type] = k_detail
}

func (a *API) Get(k_type string) string {
	k_detail, ok := a.KEY[k_type]
	if ok {
		return k_detail
	} else {
		return ""
	}
}

func (a *API) Query(userid, uri string) (interface{}, error) {
	uri_type, uri_content := ParseURI(uri)
	api_key, _ := a.KEY[uri_type]

	switch uri_type {
	case "TURING":
		return QueryTURING(userid, uri_content, api_key)
	case "QDT":
		return QueryQDT(uri_content)
	case "ZYQD_OPEN":
		return OpenZYQD(uri_content)
	case "ZYQD_CLOSE":
		return CloseZYQD(uri_content)
	default:
		return "", errors.New("Unknown type:" + uri_type)
	}

	return "", nil
}

func ParseURI(uri string) (uri_type, uri_content string) {
	uri_array := strings.SplitN(uri, `://`, 2)
	if len(uri_array) == 2 {
		return uri_array[0], uri_array[1]
	} else {
		return "", ""
	}

}

func QueryAPI(uri string) (string, error) {
	uri_type, uri_content := ParseURI(uri)
	if uri_type == "QDT" {
		return QueryQDT(uri_content)
	} else {

	}
	return "查询[" + uri_type + "]失败", errors.New("未知类型查询:" + uri_content)
}
