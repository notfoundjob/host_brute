package main

import (
	"kw0ng.top/hosts/common"
	"kw0ng.top/hosts/runner"
)

func main() {
	var hostsOptions common.HostsOptions
	common.Flag(&hostsOptions)
	common.Parse(&hostsOptions)
	runner.Run(&hostsOptions)
}
