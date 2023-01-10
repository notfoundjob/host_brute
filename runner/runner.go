package runner

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kw0ng.top/hosts/common"
)

func CheckHost(hostArrayAll []string, client *http.Client, options *common.HostsOptions) []string {
	hostSem := make(chan struct{}, options.Thread)
	var hostArrayAllUrl []string
	for _, host := range hostArrayAll {
		hostArrayAllUrl = append(hostArrayAllUrl, "http://"+host)
		hostArrayAllUrl = append(hostArrayAllUrl, "https://"+host)
	}
	hostCh := make(chan common.UrlResult, len(hostArrayAllUrl))
	for _, targetUrl := range hostArrayAllUrl {
		hostSem <- struct{}{}
		go func(targetUrl string) {
			defer func() { <-hostSem }()
			common.CheckHost(client, targetUrl, hostCh, options.Header)
		}(targetUrl)
	}
	results := make([]common.UrlResult, len(hostArrayAllUrl))
	for i := 0; i < len(hostArrayAllUrl); i++ {
		results[i] = <-hostCh
	}
	hostArray := []string{}
	for _, item := range results {
		hostArray = append(hostArray, item.Host)
	}
	close(hostCh)
	close(hostSem)
	return common.RemoveRepByMap(hostArray)
}

func Run(options *common.HostsOptions) {
	// 尝试是否能访问
	urlArray := common.ReadFile(options.Urls)
	hostArrayAll := common.ReadFile(options.Hosts)
	whiteStatusCode := make(map[int]struct{})
	if options.StatusCode != "" {
		statusCodeArray := strings.Split(options.StatusCode, ",")
		for _, value := range statusCodeArray {
			temp, _ := strconv.Atoi(value)
			whiteStatusCode[temp] = struct{}{}
		}
	}
	// proxyUrl, _ := url.Parse("http://127.0.0.1:8080")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Proxy:           http.ProxyURL(proxyUrl),
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(options.TimeOut) * time.Second,
		// 禁止重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	hostArray := CheckHost(hostArrayAll, client, options)
	sem := make(chan struct{}, options.Thread)
	hostWeightMap := []map[string]interface{}{}
	for _, targetUrl := range urlArray {
		targetUrlData := common.Check(client, targetUrl, options.Header)
		if targetUrlData.StatusCode != 0 {
			// 存储结果
			ch := make(chan common.UrlResult, len(hostArray))
			for _, host := range hostArray {
				sem <- struct{}{}
				go func(host string) {
					defer func() { <-sem }()
					common.Get(client, targetUrl, host, ch, options.Header, whiteStatusCode)
				}(host)
			}
			results := make([]common.UrlResult, len(hostArray))
			for i := 0; i < len(hostArray); i++ {
				results[i] = <-ch
			}
			resultsWeight := map[string]int{}
			hostTitleCount := 0
			reponseMd5Count := 0
			fmt.Printf("[*] URL: %s Title: %s StatusCode: %d ResponseMd5: %s\n", targetUrlData.Url, targetUrlData.Title, targetUrlData.StatusCode, targetUrlData.ResponseMd5)
			for _, item := range results {
				if item.StatusCode != 0 {
					if item.Title != targetUrlData.Title {
						fmt.Printf("[+] Title 差异 Host: %s 原始Title: %s | HostTitle: %s\n", item.Host, targetUrlData.Title, item.Title)
						resultsWeight[item.Host] += 1
						hostTitleCount += 1
					}
					if item.ResponseMd5 != targetUrlData.ResponseMd5 {
						fmt.Printf("[+] Md5   差异 Host: %s 原始Md5: %s | HostMd5: %s\n", item.Host, targetUrlData.ResponseMd5, item.ResponseMd5)
						resultsWeight[item.Host] += 1
						reponseMd5Count += 1
					}
					if item.StatusCode != targetUrlData.StatusCode {
						fmt.Printf("[+] 状态码差异 Host: %s 原始状态码: %d | Host状态码: %d\n", item.Host, targetUrlData.StatusCode, item.StatusCode)
						resultsWeight[item.Host] += 1
					}
				}
			}
			// 脏数据 md5超过10 标题超过10
			if hostTitleCount <= 10 && reponseMd5Count <= 10 {
				for host, weight := range resultsWeight {
					hostWeightMap = append(hostWeightMap, map[string]interface{}{"url": targetUrl, "host": host, "weight": weight})
				}
			}
			close(ch)
		}
	}
	close(sem)
	if len(hostWeightMap) > 0 {
		var fileName string
		if options.Output == "json" {
			fileName = common.CreateOutput("json")
			resultsJson, _ := json.Marshal(&hostWeightMap)
			common.WriteJson(fileName, string(resultsJson))
		} else if options.Output == "csv" {
			fileName = common.CreateOutput("csv")
			common.WriteCsv(fileName, hostWeightMap)
		}
	} else {
		fmt.Println("[-] 去除脏数据后未发现存在Host碰撞的Url")
	}
}

// 检测的太慢
