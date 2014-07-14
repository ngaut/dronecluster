package main

import (
	"flag"

	"github.com/ngaut/dronecluster/service"
)

var (
	server       = flag.String("server", "http://localhost:9090", "tyrant master address")
	executorUrls = flag.String("exeurls", "http://localhost:80/dronecluster.tar.gz", "executor urls")
)

func main() {
	flag.Parse()
	s := &service.ApiServer{Server: *server, ExecutorUrls: *executorUrls}
	s.Start()
}
