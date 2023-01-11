package runner

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/notfoundjob/host_brute/common"
)

func CheckHost(hostArrayAll []string, client *http.Client, options *common.HostsOptions, userAgentList []string) []string {
	hostSem := make(chan struct{}, options.Thread)
	var hostArrayAllUrl []string
	for _, host := range hostArrayAll {
		fmt.Printf("检测Host: %s\n", host)
		hostArrayAllUrl = append(hostArrayAllUrl, "http://"+host)
		hostArrayAllUrl = append(hostArrayAllUrl, "https://"+host)
	}
	hostCh := make(chan common.UrlResult, len(hostArrayAllUrl))
	rand.Seed(time.Now().UnixNano())
	for _, targetUrl := range hostArrayAllUrl {
		randomAgent := ""
		if options.RandomAgent {
			i := rand.Intn(len(userAgentList))
			randomAgent = userAgentList[i]
		}
		hostSem <- struct{}{}
		go func(targetUrl string) {
			defer func() { <-hostSem }()
			common.CheckHost(client, targetUrl, hostCh, options.Header, randomAgent)
		}(targetUrl)
	}
	results := make([]common.UrlResult, len(hostArrayAllUrl))
	for i := 0; i < len(hostArrayAllUrl); i++ {
		results[i] = <-hostCh
	}
	hostArray := []string{}
	for _, item := range results {
		if item.StatusCode == 0 {
			hostArray = append(hostArray, item.Host)
		}
	}
	close(hostCh)
	close(hostSem)
	return common.RemoveRepByMap(hostArray)
}

func Run(options *common.HostsOptions) {
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
	rand.Seed(time.Now().UnixNano())
	var userAgentList []string
	if options.RandomAgent {
		userAgentList = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5163.147 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:108.0) Gecko/20100101 Firefox/108.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
			"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1474.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.23",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		}
	}
	hostArray := CheckHost(hostArrayAll, client, options, userAgentList)
	sem := make(chan struct{}, options.Thread)
	hostWeightMap := []map[string]interface{}{}
	common.BarClass.NewOption(0, int64(len(urlArray)))
	for num, targetUrl := range urlArray {
		randomAgent := ""
		if options.RandomAgent {
			i := rand.Intn(len(userAgentList))
			randomAgent = userAgentList[i]
		}
		targetUrlData := common.Check(client, targetUrl, options.Header, randomAgent)
		if targetUrlData.StatusCode != 0 {
			common.Log(fmt.Sprintf("[*] URL: %s Title: %s StatusCode: %d ResponseMd5: %s", targetUrlData.Url, targetUrlData.Title, targetUrlData.StatusCode, targetUrlData.ResponseMd5))
			// 存储结果
			ch := make(chan common.UrlResult, len(hostArray))
			for _, host := range hostArray {
				randomAgent := ""
				if options.RandomAgent {
					i := rand.Intn(len(userAgentList))
					randomAgent = userAgentList[i]
				}
				sem <- struct{}{}
				go func(host string) {
					defer func() { <-sem }()
					common.Get(client, targetUrl, host, ch, options.Header, whiteStatusCode, randomAgent)
				}(host)
			}
			results := make([]common.UrlResult, len(hostArray))
			for i := 0; i < len(hostArray); i++ {
				results[i] = <-ch
			}
			resultsWeight := map[string]int{}
			resultsTitle := map[string]string{}
			resultsStatusCode := map[string]int{}
			hostTitleCount := 0
			reponseMd5Count := 0
			for _, item := range results {
				if item.StatusCode != 0 {
					if item.Title != targetUrlData.Title {
						common.Log(fmt.Sprintf("[+] Title 差异 Host: %s 原始Title: %s | HostTitle: %s", item.Host, targetUrlData.Title, item.Title))
						resultsWeight[item.Host] += 1
						resultsTitle[item.Host] = item.Title
						hostTitleCount += 1
					}
					if item.ResponseMd5 != targetUrlData.ResponseMd5 {
						common.Log(fmt.Sprintf("[+] Md5   差异 Host: %s 原始Md5: %s | HostMd5: %s", item.Host, targetUrlData.ResponseMd5, item.ResponseMd5))
						resultsWeight[item.Host] += 1
						reponseMd5Count += 1
					}
					if item.StatusCode != targetUrlData.StatusCode {
						common.Log(fmt.Sprintf("[+] 状态码差异 Host: %s 原始状态码: %d | Host状态码: %d", item.Host, targetUrlData.StatusCode, item.StatusCode))
						resultsWeight[item.Host] += 1
						resultsStatusCode[item.Host] = item.StatusCode
					}
				}
			}
			if hostTitleCount <= options.Max && reponseMd5Count <= options.Max {
				for host, weight := range resultsWeight {
					statusCode := resultsStatusCode[host]
					if statusCode == 0 {
						statusCode = targetUrlData.StatusCode
					}
					hostWeightMap = append(hostWeightMap, map[string]interface{}{
						"url":        targetUrl,
						"host":       host,
						"weight":     weight,
						"title":      resultsTitle[host],
						"statuscode": statusCode,
					})
				}
			}
			close(ch)
		}
		common.BarInt = int64(num + 1)
	}
	common.BarClass.Play(int64(len(urlArray)))
	fmt.Println()
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
		fmt.Println("[-] 去除脏数据后, 未发现存在Host碰撞的Url!")
	}
}
