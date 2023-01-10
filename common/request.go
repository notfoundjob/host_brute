package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// func GetReq(client *resty.Client, url, host string) (string, int) {
// 	if host != "" {
// 		log.Printf(host)
// 		client.SetHeader("Host", host)
// 	}
// 	resp, err := client.R().Get(url)
// 	if err != nil {
// 		log.Printf("请求发生错误: %s", err)
// 		return "", 0
// 	}
// 	return string(resp.Body()), resp.StatusCode()
// }
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

func Check(client *http.Client, targetUrl string, headers []string) UrlResult {
	log.Printf("检测: %s", targetUrl)
	title, responseMd5, statusCode := "", "", 0
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Printf("NewRequest发生错误: %s", err)
		return UrlResult{
			StatusCode: 0,
		}
	}
	req.Header.Add("User-Agent", userAgent)
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
		log.Printf("请求发生错误: %s", err)
		return UrlResult{
			StatusCode: 0,
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取返回内容发生错误: %s", err)
		statusCode = 0
	} else {
		title = GetTitle(string(body))
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

func Get(client *http.Client, targetUrl, host string, ch chan<- UrlResult, headers []string, whiteStatusCode map[int]struct{}) {
	log.Printf("访问: %s Host: %s", targetUrl, host)
	title, responseMd5, statusCode := "", "", 0
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Printf("NewRequest发生错误: %s", err)
	}
	req.Host = host
	req.Header.Add("User-Agent", userAgent)
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
		log.Printf("请求发生错误: %s", err)
		ch <- UrlResult{
			StatusCode: 0,
		}
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("读取返回内容发生错误: %s", err)
		} else {
			title = GetTitle(string(body))
			statusCode = resp.StatusCode
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
				log.Printf("丢弃: %s Host: %s 状态码: %d", targetUrl, host, statusCode)
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
		// time.Sleep(5 * time.Second)
	}
}

func CheckHost(client *http.Client, targetUrl string, ch chan<- UrlResult, headers []string) {
	log.Printf("检测: %s", targetUrl)
	statusCode := 0
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Printf("NewRequest发生错误: %s", err)
	}
	req.Header.Add("User-Agent", userAgent)
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
		// log.Printf("请求发生错误: %s", err)
		ch <- UrlResult{
			Host:       req.Host,
			StatusCode: statusCode,
		}
	} else {
		// time.Sleep(5 * time.Second)
		ch <- UrlResult{
			StatusCode: statusCode,
		}
	}
}
