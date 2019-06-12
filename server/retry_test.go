package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestRetry(t *testing.T) {
	go withRetryHttp("10086")
	values := url.Values{}
	values.Add("addr", "")
	resp, err := http.PostForm("http://127.0.0.1:10086/register", values)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(body))
}
