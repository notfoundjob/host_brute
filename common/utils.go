package common

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-basic/uuid"
)

func GetMd5(text []byte) string {
	srcCode := md5.Sum(text)
	code := fmt.Sprintf("%x", srcCode)
	return string(code)
}

func GetTitle(text string) string {
	re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	matches := re.FindStringSubmatch(string(text))
	// fmt.Println(matches) 全量,括号
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func ReadFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open %s error, %v\n", filename, err)
		os.Exit(0)
	}
	defer file.Close()
	var content []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			content = append(content, scanner.Text())
		}
	}
	return content
}

func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l && e != "" {
			result = append(result, e)
		}
	}
	return result
}

func CreateOutput(fileType string) string {
	var fileName string
	if fileType == "json" {
		fileName = "result.json"
		_, outPathErr := os.Stat(fileName)
		if os.IsNotExist(outPathErr) {
			os.Create(fileName)
			fmt.Println("[+] 输出为: ", fileName)
		} else {
			fileName = "result-" + uuid.New() + ".json"
			fmt.Println("[+] result.json 文件已存在, 输出为: ", fileName)
			os.Create(fileName)
		}
	} else if fileType == "csv" {
		fileName = "result.csv"
		_, outPathErr := os.Stat(fileName)
		if os.IsNotExist(outPathErr) {
			os.Create(fileName)
			fmt.Println("[+] 输出为: ", fileName)
		} else {
			fileName = "result-" + uuid.New() + ".csv"
			fmt.Println("[+] result.csv 文件已存在, 输出为: ", fileName)
			os.Create(fileName)
		}
		f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		io.WriteString(f, "Url,Host,Weight,Title,StatusCode\n")
		defer f.Close()
	}
	return fileName
}

func Log(content string) {
	if TerminalWidth > len(content) {
		fmt.Printf("\r" + content + strings.Repeat(" ", TerminalWidth-len(content)) + "\n")
	} else {
		fmt.Printf("\r" + content + "\n")
	}
	BarClass.Play(BarInt)
}

func WriteJson(fileName, result string) {
	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	io.WriteString(f, result+"\n")
	defer f.Close()
}

func WriteCsv(fileName string, result []map[string]interface{}) {
	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	for _, item := range result {
		io.WriteString(f, item["url"].(string)+","+item["host"].(string)+","+strconv.Itoa(item["weight"].(int))+","+item["title"].(string)+","+strconv.Itoa(item["statuscode"].(int))+"\n")
	}
	defer f.Close()
}
