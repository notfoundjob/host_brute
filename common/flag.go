package common

import (
	"flag"
	"fmt"
)

func (i *ArrayFlags) String() string {
	return fmt.Sprint(*i)
}
func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func Flag(Info *HostsOptions) {
	flag.StringVar(&Info.Hosts, "hf", "", "Host文件")
	flag.StringVar(&Info.Urls, "uf", "", "Url文件")
	flag.BoolVar(&Info.RandomAgent, "random-agent", false, "开启随机User-Agent")
	flag.Var(&Info.Header, "H", "请求Header 支持解析多个 -H")
	flag.IntVar(&Info.Thread, "t", 32, "最大线程数量")
	flag.IntVar(&Info.Max, "m", 50, "单个Url如果出现N次Host 则认定为脏数据")
	flag.IntVar(&Info.TimeOut, "timeout", 10, "请求超时时间(秒)")
	flag.StringVar(&Info.StatusCode, "mc", "", "状态码 200,500")
	flag.StringVar(&Info.Output, "o", "csv", "输出文件格式 json 或 csv")
	flag.Parse()
}
