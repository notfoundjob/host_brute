package main

import (
	"github.com/notfoundjob/host_brute/common"
	"github.com/notfoundjob/host_brute/runner"
)

func main() {
	var hostsOptions common.HostsOptions
	common.Flag(&hostsOptions)
	common.Parse(&hostsOptions)
	runner.Run(&hostsOptions)
}
