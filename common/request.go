package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/axgle/mahonia"
	"golang.org/x/net/html/charset"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

func Check(client *http.Client, targetUrl string, headers []string, randomAgent string) UrlResult {
	Log("检测: " + targetUrl)
	title, responseMd5, statusCode := "", "", 0
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		Log(fmt.Sprintf("NewRequest发生错误: %s", err))
		return UrlResult{
			StatusCode: 0,
		}
	}
	if randomAgent != "" {
		req.Header.Add("User-Agent", randomAgent)
	} else {
		req.Header.Add("User-Agent", userAgent)
	}
	if len(headers) > 0 {
		for _, item := range headers {
			temp := strings.Split(item, ":")
			req.Header.Add(temp[0], temp[1])
		}
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		Log(fmt.Sprintf("请求发生错误: %s", err))
		return UrlResult{
			StatusCode: 0,
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log(fmt.Sprintf("读取返回内容发生错误: %s", err))
		statusCode = 0
	} else {
		_, encoding, _ := charset.DetermineEncoding(body, resp.Header.Get("Content-Type"))
		if encoding == "gbk" {
			bodyStr := mahonia.NewDecoder("gbk").ConvertString(string(body))
			title = GetTitle(bodyStr)
		} else {
			title = GetTitle(string(body))
		}
		statusCode = resp.StatusCode
		responseMd5 = GetMd5(body)
	}
	return UrlResult{
		Url:         targetUrl,
		Title:       title,
		StatusCode:  statusCode,
		ResponseMd5: responseMd5,
	}
}

func Get(client *http.Client, targetUrl, host string, ch chan<- UrlResult, headers []string, whiteStatusCode map[int]struct{}, randomAgent string) {
	Log(fmt.Sprintf("测试: %s Host: %s", targetUrl, host))
	title, responseMd5, statusCode := "", "", 0
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		Log(fmt.Sprintf("NewRequest发生错误: %s", err))
	}
	req.Host = host
	if randomAgent != "" {
		req.Header.Add("User-Agent", randomAgent)
	} else {
		req.Header.Add("User-Agent", userAgent)
	}
	if len(headers) > 0 {
		for _, item := range headers {
			temp := strings.Split(item, ":")
			req.Header.Add(temp[0], temp[1])
		}
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		Log(fmt.Sprintf("请求发生错误: %s", err))
		ch <- UrlResult{
			StatusCode: 0,
		}
	} else {
		statusCode = resp.StatusCode
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Log(fmt.Sprintf("读取返回内容发生错误: %s", err))
		} else {
			_, encoding, _ := charset.DetermineEncoding(body, resp.Header.Get("Content-Type"))
			if encoding == "gbk" {
				bodyStr := mahonia.NewDecoder("gbk").ConvertString(string(body))
				title = GetTitle(bodyStr)
			} else {
				title = GetTitle(string(body))
			}
			responseMd5 = GetMd5(body)
		}
		if len(whiteStatusCode) != 0 {
			// 如果用户传入了状态码白名单
			if _, ok := whiteStatusCode[statusCode]; ok {
				// 判断当前状态码是否在白名单内
				ch <- UrlResult{
					Url:         targetUrl,
					Host:        host,
					Title:       title,
					StatusCode:  statusCode,
					ResponseMd5: responseMd5,
				}
			} else {
				Log(fmt.Sprintf("丢弃: %s Host: %s 状态码: %d", targetUrl, host, statusCode))
				ch <- UrlResult{
					StatusCode: 0,
				}
			}
		} else {
			ch <- UrlResult{
				Url:         targetUrl,
				Host:        host,
				Title:       title,
				StatusCode:  statusCode,
				ResponseMd5: responseMd5,
			}
		}
	}
}

func CheckHost(client *http.Client, targetUrl string, ch chan<- UrlResult, headers []string, randomAgent string) {
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		Log(fmt.Sprintf("NewRequest发生错误: %s", err))
	}
	if randomAgent != "" {
		req.Header.Add("User-Agent", randomAgent)
	} else {
		req.Header.Add("User-Agent", userAgent)
	}
	if len(headers) > 0 {
		for _, item := range headers {
			temp := strings.Split(item, ":")
			req.Header.Add(temp[0], temp[1])
		}
	}
	resp, doErr := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if doErr != nil {
		// 访问失败 保留
		ch <- UrlResult{
			Host:       req.Host,
			StatusCode: 0,
		}
	} else {
		// 访问成功 抛弃
		ch <- UrlResult{
			StatusCode: 200,
		}
	}
}
