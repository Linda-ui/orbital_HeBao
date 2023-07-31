package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	utils "github.com/Linda-ui/orbital_HeBao/hertz_gateway/utils"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

var (
	idlPath       string
	svcMethodName string
	payload       string
	port          string
	concurrent    int
	total         int64
	bufsize       int
	sleepTime     int
)

func initFlags() {
	flag.StringVar(&idlPath, "idlpath", utils.GetProjectIDLRoot()+"/echo.thrift", "the path to the service IDL file")
	flag.StringVar(&svcMethodName, "method", "EchoMethod", "the name of the service method to be called")
	flag.StringVar(&payload, "d", "{\"msg\":\"goodbye\"}", "the data to be sent to the server")
	flag.StringVar(&port, "p", "8870", "the ports to be used for the benchmarking")
	flag.IntVar(&concurrent, "c", 100, "the number of concurrent requests")
	flag.Int64Var(&total, "n", 200000, "the total number of requests")
	flag.IntVar(&bufsize, "b", 1024, "the size of the buffer for each request")
	flag.IntVar(&sleepTime, "s", 0, "the time to sleep between each request")
	flag.Parse()
}

// a general RPC benchmarking function that benches the performance of calling a server directly without going through the API gateway.
func main() {
	initFlags()
	svcName := utils.ExtractServiceName(idlPath)

	// creating a new generic client.
	p, err := generic.NewThriftFileProvider(idlPath)
	if err != nil {
		log.Fatalf("creating new thrift file provider for %v service failed: %v", svcName, err)
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		log.Fatalf("creating new generic instance for %v service failed: %v", svcName, err)
	}

	cli, err := genericclient.NewClient(svcName, g)
	if err != nil {
		log.Fatalf("creating new generic client for %v service failed: %v", svcName, err)
	}

	r := runner.NewRunner()
	recorder := perf.NewRecorder(fmt.Sprintf("%s@Client", svcName))

	handler := func() error {
		_, err := cli.GenericCall(context.Background(), svcMethodName, payload, callopt.WithHostPort("127.0.0.1:"+port))
		return err
	}

	recorder.Begin()
	r.Run(svcName+" @port"+port, handler, concurrent, total, bufsize, sleepTime)
	recorder.End()
	recorder.Report()
}
