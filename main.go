/**
 * main process
 * Usage: wxpay-gateway[ -v]
 * Rosbit Xu
 */
package main

import (
	"os"
	"fmt"
	"github.com/rosbit/go-wxpay-gateway/conf"
)

// variables set via go build -ldflags
var (
	buildTime string
	osInfo    string
	goInfo    string
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		ShowInfo("name",       os.Args[0])
		ShowInfo("build time", buildTime)
		ShowInfo("os name",    osInfo)
		ShowInfo("compiler",   goInfo)
		return
	}

	if err := conf.CheckGlobalConf(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(3)
		return
	}
	conf.DumpConf()

	if err := StartService(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(4)
	}
	os.Exit(0)
}

func ShowInfo(prompt, info string) {
	if info != "" {
		fmt.Printf("%10s: %s\n", prompt, info)
	}
}
