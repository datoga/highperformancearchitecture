package main

import "flag"

import _ "github.com/datoga/highperformancearchitecture/cmd/common/flags"

var pace = flag.String("pace", "1s", "pace between random messages")

func init() {
	flag.Parse()
}
