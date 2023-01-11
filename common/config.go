package common

type ArrayFlags []string

type HostsOptions struct {
	Hosts       string
	Urls        string
	RandomAgent bool
	Header      ArrayFlags
	Thread      int
	Max         int
	TimeOut     int
	StatusCode  string
	Output      string
}

type UrlResult struct {
	Url         string
	Host        string
	Title       string
	StatusCode  int
	ResponseMd5 string
}

type HostWeight struct {
	Host   string
	Weight int
}

type Bar struct {
	total int64
}

var TerminalWidth int

var BarInt int64

var BarClass Bar
