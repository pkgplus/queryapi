package queryapi

import (
	"encoding/json"
)

const (
	URL_ZYQD_OPEN = "http://127.0.0.1/zyqd/open"
	URL_ZYQD_CODE = "http://127.0.0.1/zyqd/close"
)

type QDResPonse struct {
	Ret    string
	ErrMsg string
}

func OpenZYQD(ip string) ([]byte, error) {
	exe_client := &MyHttpClient{
		Method:         "GET",
		Url:            URL_ZYQD_OPEN + "?ip=" + ip,
		SaveCookieFlag: false,
		ContentType:    `text/html`,
	}
	err := exe_client.Do()
	if err != nil {
		return exe_client.ContentBytes, err
	} else {
		return exe_client.ContentBytes, nil
	}
}

func CloseZYQD(code string) (string, error) {
	exe_client := &MyHttpClient{
		Method:         "GET",
		Url:            URL_ZYQD_CODE + "?code=" + code,
		SaveCookieFlag: false,
		ContentType:    `text/html`,
	}
	err := exe_client.Do()
	if err != nil {
		return "", err
	}

	ret := &QDResPonse{}
	json_err := json.Unmarshal(exe_client.ContentBytes, ret)
	if json_err != nil {
		return "", json_err
	} else {
		return ret.ErrMsg, nil
	}
}
