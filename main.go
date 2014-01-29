package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
	Bbs "github.com/onsi-experimental/runtime-schema/bbs"
	"github.com/vito/gordon"

	"github.com/onsi-experimental/executor/executor"
)

var wardenNetwork = flag.String(
	"wardenNetwork",
	"unix",
	"network mode for warden server (tcp, unix)",
)

var wardenAddr = flag.String(
	"wardenAddr",
	"/tmp/warden.sock",
	"network address for warden server",
)

var etcdMachines = flag.String(
	"etcdMachines",
	"http://127.0.0.1:4001",
	"comma-separated list of etcd addresses (http://ip:port)",
)

func main() {
	flag.Parse()

	etcdAdapter := etcdstoreadapter.NewETCDStoreAdapter(
		strings.Split(*etcdMachines, ","),
		workerpool.NewWorkerPool(10),
	)

	bbs := Bbs.New(etcdAdapter)
	err := etcdAdapter.Connect()
	if err != nil {
		log.Fatalln("failed to get etcdAdapter to connect")
	}

	wardenClient := gordon.NewClient(&gordon.ConnectionInfo{
		Network: *wardenNetwork,
		Addr:    *wardenAddr,
	})

	err = wardenClient.Connect()
	if err != nil {
		log.Fatalln("warden is not up!", err)
	}

	executor := executor.New(bbs, wardenClient)

	executor.HandleRunOnces()

	fmt.Println("Watching for RunOnces!")
	select {}
}
