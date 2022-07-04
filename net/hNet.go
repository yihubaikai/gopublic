package hNet

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	//_ "github.com/go-sql-driver/mysql"
)

//url 编码  "net/url"
func Urlencode(val string) string {
	var str string
	u := url.Values{}
	u.Set("", val)
	str = u.Encode()
	str = str[1:]
	return str
}

//url解码  "net/url"
func Urldecode(val string) string {
	urlDecodeStr, _ := url.QueryUnescape(val)
	return urlDecodeStr
}

//	encdatas := `00000000000000000000000000000000,money,5000.99` //{"para_id":"gameb","order_no":"00000000000000000000000000000000","money":"5000.99"}
//  fmt.Println("加密字符串:" + string(encdatas))
func Myencode(p []byte) string {
	encodedStr := hex.EncodeToString([]byte(p))
	return encodedStr
}

//encodedStr := myencode([]byte(encdatas))
//fmt.Println(encodedStr)
func Mydecode(p string) string {
	decodeStr, _ := hex.DecodeString(p)
	return string(decodeStr)
}

//get请求
//调用demo
//	s := make(map[string]string)
//	s["wd"] = "牛魔王之红孩儿诞生"
//	s["act"] = "我是get请求"
//	r := httpget("http://pay.ggpaygg.com/debug/recvtest.php", s)
//	fmt.Println(r)

func Httpgetz(desurl string, para_kv map[string]string) string {
	fullurl := desurl
	if para_kv != nil {
		u := url.Values{}
		for k, v := range para_kv {
			//fmt.Println(k, "*", v)
			u.Set(k, v)
		}
		fullurl = desurl + "?" + u.Encode()
	}

	resp, err := http.Get(fullurl)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close() //not ok
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

func Httpget(desurl string, para_kv map[string]string) string {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		err := recover()
		if err != nil {
			fmt.Println("拦截:", err) // 这里的err其实就是panic传入的内容，55
			//return ""
		}
	}()
	return Httpgetz(desurl, para_kv)
}

//post请求
//	s := make(map[string]string)
//	s["wd"] = "牛魔王之红孩儿诞生"
//  s["act"] = "我是post请求"
//	p := httppost("http://pay.ggpaygg.com/debug/recvtest.php", s)
//	fmt.Println(p)
func Httppostz(desurl string, para_kv map[string]string) string {
	u := url.Values{}
	for k, v := range para_kv {
		//fmt.Println(k, "*", v)
		u.Set(k, v)
	}
	fullurl := u.Encode()
	//fmt.Println(paramstr)

	resp, err := http.Post(desurl, "application/x-www-form-urlencoded;charset=utf-8", strings.NewReader(fullurl))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close() //not ok
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}
func HttppostJson(desurl string, para_kv map[string]string) string {
	u := url.Values{}
	for k, v := range para_kv {
		//fmt.Println(k, "*", v)
		u.Set(k, v)
	}
	fullurl := u.Encode()
	//fmt.Println(paramstr)

	resp, err := http.Post(desurl, "application/json;charset=utf-8", strings.NewReader(fullurl))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close() //not ok
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}
func Httppost(desurl string, para_kv map[string]string) string {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		err := recover()
		if err != nil {
			fmt.Println("拦截:", err) // 这里的err其实就是panic传入的内容，55
			//return ""
		}
	}()
	return Httppostz(desurl, para_kv)
}
