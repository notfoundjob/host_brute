package common

type ArrayFlags []string

type HostsOptions struct {
	Hosts      string
	Urls       string
	UserAgent  bool
	Header     ArrayFlags
	Thread     int
	Proxy      string
	TimeOut    int
	StatusCode string
	Output     string
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

// var myFlags arrayFlags
