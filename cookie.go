package queryapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func AddCookies(req *http.Request, cookie_file string) error {
	fi, err := os.Open(cookie_file)
	if err != nil {
		return err
	}
	defer fi.Close()

	f_bytes, read_err := ioutil.ReadAll(fi)
	if read_err != nil {
		return read_err
	}

	cookies := make([]*http.Cookie, 0)
	json_err := json.Unmarshal(f_bytes, &cookies)
	if json_err != nil {
		return json_err
	}

	for _, cookie := range cookies {
		//fmt.Println(cookie)
		req.AddCookie(cookie)
	}

	return nil
}

func SaveCookies(cookies []*http.Cookie, cookie_file string) error {
	cookie_bytes, cookie_err := json.Marshal(cookies)

	if cookie_err != nil {
		return cookie_err
	}
	cookie_w_err := ioutil.WriteFile(cookie_file, cookie_bytes, 0666)
	if cookie_w_err != nil {
		return cookie_w_err
	}

	return nil
}
