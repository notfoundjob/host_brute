# host_brute

找不到工作的两人组😢一位思考者🤔+一位编写者✍️=有想法的项目

相关思考：[黑客的睡前一思](https://mp.weixin.qq.com/s/QCPUwgwhnDtuuY656Ec6XQ)

希望大家多提issue，让我们更好的为大家服务，respect

> 资产信息扩大器，通过host碰撞发现企业隐形资产，帮助攻击者发现更多攻击面。

## 功能

- [x] url和host存活检测

- [x] 标题差异检测

- [x] 状态码差异检测

- [x] 页面内容差异检测
- [x] 自定义状态码（白名单）

## 使用

```
host_brute -h
```

支持如下参数

```
Usage:
  -H value
    	请求Header 支持解析多个 -H
  -hf string
    	Host文件
  -mc string
    	状态码 200,500
  -o string
    	输出文件格式 json 或 csv (default "csv")
  -random-agent
    	开启随机User-Agent (未支持)
  -t int
    	最大线程数量 (default 32)
  -timeout int
    	请求超时时间(秒) (default 10)
  -uf string
    	Url文件
```

## 举例

url文件

```
http://1.2.3.4:8090
http://1.2.3.4
https://1.2.3.4:8765
```

host文件

```
inner.example.com
vir.example.com
local.example.local
localhost
test
mailtest
```

命令

```
host_brute -uf url.txt -hf host.txt -t 50 -timeout 5 -mc 200,301,302
```

## Tips

权重越高，准确率越高

## 遇到的问题

- 页面内容随机变化

## 未来

- [ ] 进度条
- [ ] 随机UA
- [ ] 页面相似度
