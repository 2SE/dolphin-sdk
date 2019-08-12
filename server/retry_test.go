package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestRetry(t *testing.T) {
	//go withRetryHttp("10086")
	values := url.Values{}
	values.Add("addr", "http://192.168.10.169:9527")
	resp, err := http.PostForm("http://192.168.9.101:9797/register", values)
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
	//values.Add("addr", "http://192.168.10.169:9527")
	resp, err = http.PostForm("http://192.168.9.101:9999/register", values)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	resp, err = http.PostForm("http://192.168.10.169:10003/register", values)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(body))
}

func TestReg(t *testing.T) {
	reg, _ = regexp.Compile(`^[a-zA-Z0-9]{6,20}$`)

	t.Log(reg.MatchString("asd123"))
}
