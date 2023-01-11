package common

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func Parse(options *HostsOptions) {
	_, hostsPathErr := os.Stat(options.Hosts)
	if options.Hosts == "" || os.IsNotExist(hostsPathErr) {
		fmt.Println("Host文件不存在")
	}
	_, urlsPathErr := os.Stat(options.Urls)
	if options.Urls == "" || os.IsNotExist(urlsPathErr) {
		fmt.Println("Url文件不存在")
		os.Exit(0)
	}
	TerminalWidth, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
}
