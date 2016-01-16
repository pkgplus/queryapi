package queryapi

import (
	"code.google.com/p/mahonia"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type MyHttpClient struct {
	Method         string
	Url            string
	Decode         string
	CookieFile     string
	SaveCookieFlag bool
	Refer          string
	ContentType    string
	ContentBytes   []byte
	PostData       *url.Values
}

func (m_client *MyHttpClient) Do() error {
	m_client.ContentBytes = nil

	//New Request
	var req *http.Request
	var req_err error
	if m_client.PostData != nil {
		req, req_err = http.NewRequest(m_client.Method, m_client.Url, strings.NewReader(m_client.PostData.Encode()))
	} else {
		req, req_err = http.NewRequest(m_client.Method, m_client.Url, nil)
	}

	log.Println(req.URL.String())
	if req_err != nil {
		return req_err
	}

	//add header
	if m_client.ContentType == "" {
		m_client.ContentType = "text/html;charset=gb2312"
		m_client.Decode = "gb2312"
	}
	req.Header.Set("Content-Type", m_client.ContentType)
	req.Header.Set("Referer", m_client.Refer)

	if m_client.CookieFile != "" {
		AddCookies(req, m_client.CookieFile)
	}

	//New HTTP CLIENT
	client := &http.Client{}
	resp, client_err := client.Do(req)
	if client_err != nil {
		return client_err
	}
	defer resp.Body.Close()

	//cookie
	if m_client.SaveCookieFlag && len(resp.Cookies()) > 0 {
		SaveCookies(resp.Cookies(), m_client.CookieFile)
	}

	//login result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return read_err
	}

	//convert gb2312 to utf-8
	if m_client.Decode != "" {
		dec := mahonia.NewDecoder(m_client.Decode)
		ret, ok := dec.ConvertStringOK(string(resp_body))
		if !ok {
			return errors.New("convert gb2312 to utf-8 err!")
		}
		m_client.ContentBytes = []byte(ret)
	} else {
		m_client.ContentBytes = resp_body
	}

	return nil
}
