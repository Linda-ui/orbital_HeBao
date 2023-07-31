# API Gateway based on CloudWeGo Projects
This is the orbital23 project for team HeBao.
<br></br>


## Proposed Level of Achievement
Artemis
<br></br>


## Documentation

The detailed documentation of this project including the system design, core features, and testing details can be accessed at [this page](https://lively-cereal-3a1.notion.site/API-Gateway-based-on-CloudWeGo-Projects-e59b321209f043a0a625a8c98364c838?pvs=4).
<br></br>


## Overview

This repository contains the implementation of a fully functioning API gateway using CloudWeGo's Kitex (an RPC framework) and Hertz (an HTTP framework). 

The four core features we implemented for the gateway are:
1. **Service registration and discovery:** Nacos is integrated in our gateway as the service registry. Kitex service instances register their addresses on the Nacos server, and Kitex clients can discover these instances through the Nacos server to establish connections with the corresponding servers.

2. **JSON mapping generic call:** the gateway enables protocol translation between frontend-gateway communication (HTTP in JSON) and gateway-backend communication (Thrift binary protocol). It also allows the dynamic invocation of service methods through generic calls.

3. **Load balacing:** Kitex's default load balancer is integrated in our gateway so that loads (request traffics) can be evenly distributed on the three service instances created for each Kitex service. This feature can prevent overloading and breakdown of any single service instance during peak loads.

4. **IDL-service mapping and dynamic update:** our gateway maintains a dynamic map between each interface definition language (IDL) file and its corresponding Kitex service. When an IDL file describing a service is added/deleted, the map will update itself to reflect changes in the managed services of our gateway, even when the gateway is still running. The client would not face disruptions of the gateway service when the backend services are being added or deleted.

We have also implemented two backend Kitex test services. For each Kitex service, we created three service instances (hosted on different ports).
<br></br>



## Quick Start

First, [download the latest version of Nacos](https://github.com/alibaba/nacos/releases)

Then, run 
```shell
# run at project root directory
make gateway
```
to start the Nacos server and the API gateway.

You can also start the nacos server separately with
```shell
make nacos
```

If the gateway starts succesfully, you should see something like:
```shell
nacos is starting with standalone
nacos is starting，you can check the ~/nacos/logs/start.out
...
...
[Info] HERTZ: HTTP server listening on address=127.0.0.1:8080
```
To check if Nacos is running, check `~/nacos/logs/start.out` and click the console address it provides. You should be able to see the local Nacos console.

To test if the API gateway is running.
```shell
curl http://localhost:8080/

# should return "the api gateway is running"
```

Then, start the backend servers with
```shell
make services
```

You should be able to see echo and sum services in your Nacos console, under `ServiceManagement -> Service List`.

To send a simple request to the gateway, copy the following to your terminal:
```shell
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"msg":"hello"}' \
    -w '\n' 
```
A response will be returned:
```shell
{"err_code":0,"err_message":"success","msg":"hello"}
```

After using the gateway, run
```shell
make stop
```
to do the clean up and shut down all running servers.
<br></br>

## Unit Test

Several chosen packages are unit tested (see documentatinos for details). 

Run the following commands to run all tests and check test coverage:
```shell
go test ./...
go test ./... -cover
```
You can also run each unit test separately. Unit test files are those with file name `*_test.go`.

Note that when you run the first command above, the `TestIntegrationGateway` test may fail. This is because it is an integration test and requires starting up the gateway and backend services first (see the next section for details).
<br></br>


## Integration Test

The gateway is customised to handle business errors like invalid user requests. Depending on the user inputs, various error messages will be returned with corresponding error codes (in the form of a response).

Navigate to the `/integration_test` directory, and run the go tests in verbose mode to see all our integration testcases:
```shell
go test ./... -v
```

You can also test using cURL to view the exact response the gateway sends back for each of the request, we have provided the commands to send various valid/invalid requests in `scripts/echo_request.sh` and `scripts/sum_request.sh`. You can find the expected response in the comments above each request.
<br></br>


## Testing the Dynamic Update Feature

First, start running the gateway and the backend services, `echo` and `sum`.

### Add a service
The official documentation for creating a new Kitex service can be found [here](https://www.cloudwego.io/docs/kitex/getting-started/). You can also follow our instructions below to create a simplified version for testing. If you face any difficulties, please refer to Kitex's official documentations for debugging.

This example adds an `length` Kitex service that returns the length of the string input. You can also add another service with a different service name. 

Under the project root directory, create a new Thrift file that defines this `length` service. This is used for Kitex's code generation tool to generate skeleton code.

You MUST name your IDL file the same as your service name. In this case, name the file `length.thrift`.
```thrift
namespace go length

struct LengthReq {
  1: string msg
}

struct LengthResp {
  1: i64 strlen
}

service LengthSvc {
  LengthResp LengthMethod(1: LengthReq req)
}
```
You can also add in other requests/response structs and service methods as needed. 
After saving this file, use the Kitex compiler to compile the IDL file to generate server-side code.
```shell
kitex -module github.com/Linda-ui/orbital_HeBao -service length length.thrift
```
Note that each time you add in new content to the thrift file, you have to regenerate and update the code with the same command.

After running the above command, you will see the generated project layout:
```
.
|-- build.sh
|-- length.thrift
|-- handler
|   `-- handler.go
|-- kitex_gen
|   `-- length
|       |-- lengthsvc
|       |   |-- client.go
|       |   |-- invoker.go
|       |   |-- lengthsvc.go
|       |   `-- server.go
|       |-- k-consts.go
|       |-- k-length.go
|       `-- length.go
|-- main.go
|-- kitex_info.yaml
`-- script
    `-- bootstrap.sh
```
Notice that there is not a subdirectory called `handler` for the generated code. Please create that subdirectory yourself and put `handler.go` inside. Meanwhile, you should also change the package name of `handler.go` from `main` to `handler`. This can be achieved by changing the first line of `handler.go`.

For cleaner code structure, please place all the above generated file into a new subdirectory `length` under `/kitex_services`. Note that you may have to rewrite the import path of package `length` in all generated files with this import package. However, skipping the code reorganisation will not affect the correctness of code. 

For the ease of illustration, the code below will be based on writing this service in the root directory.

All server-side method logic should be implemented in `handler.go`. Below is an example implementation of the length service.
```go
package handler

import (
	"context"
	length "github.com/Linda-ui/orbital_HeBao/kitex_gen/length"
)

// LengthSvcImpl implements the last service interface defined in the IDL.
type LengthSvcImpl struct{}

// LengthMethod implements the LengthSvcImpl interface.
func (s *LengthSvcImpl) LengthMethod(ctx context.Context, req *length.LengthReq) (resp *length.LengthResp, err error) {
	// add this line.
	return &length.LengthResp{StrLen: int64(len(req.Msg))}, nil
}
```

Next, modify the `main.go` file as below so that it registers the service with Nacos, has a service name called "length", and uses a customised port number 8893. Avoid using ports already in use. By default, the gateway runs on port 8080, the echo service runs on ports 8870, 8871, 8872, the sum service runs on ports 9870, 9871, 9872.
```go
package main

import (
	"log"
	"net"

	length "github.com/Linda-ui/orbital_HeBao/kitex_gen/length/lengthsvc"
        handler "github.com/Linda-ui/orbital_HeBao/handler"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
)

func main() {
	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		klog.Fatal(err)
	}

	svr := length.NewServer(
		new(handler.LengthSvcImpl),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "length"}),
		server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8893}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
```
Next, run `main.go` to get the server running.
```shell
go run main.go
```
You should see something like below:
```shell
2023/07/30 16:16:59.243099 logger.go:45: [Info] udp server start, port: 55045
2023/07/30 16:16:59.243642 server.go:81: [Info] KITEX: server listen at addr=127.0.0.1:8893
2023/07/30 16:17:00.244740 logger.go:45: [Info] register instance namespaceId:<>,serviceName:<DEFAULT_GROUP@@length> with instance:<{"valid":false,"marked":false,"instanceId":"","port":8893,"ip":"127.0.0.1","weight":10,"metadata":{},"clusterName":"DEFAULT","serviceName":"","enabled":true,"healthy":true,"ephemeral":true}>
2023/07/30 16:17:00.249035 logger.go:45: [Info] adding beat: <{"ip":"127.0.0.1","port":8893,"weight":10,"serviceName":"DEFAULT_GROUP@@length","cluster":"DEFAULT","metadata":{},"scheduled":false}> to beat map
2023/07/30 16:17:00.249103 logger.go:45: [Info] namespaceId:<> sending beat to server:<{"ip":"127.0.0.1","port":8893,"weight":10,"serviceName":"DEFAULT_GROUP@@length","cluster":"DEFAULT","metadata":{},"scheduled":false}>
```
This means you have successfully run the Kitex server for the `length`service. 

In order for our gateway to recognise our newly added service, move `length.thrift` into the `/idl` directory. The addition of IDL file in the `/idl` directory tells our gateway that a new service named `length` is added.

Test this new service through the gateway by:  
```shell
curl -X POST http://localhost:8080/gateway/length/LengthMethod \
    -H "Content-Type: application/json" \
    -d '{"msg":"hello"}' \
    -w '\n'
```
You should see:
```shell
{"err_code":0,"err_message":"success","strlen":5}
```
Congratulations! You have just added a new Kitex service to the gateway while it is still running, leveraging the dynamic update feature. Note that if you want to change the server method logic, you have to shut down the running Kitex server first and restart it when you finish modifying.
<br></br>

### Delete a service
Simply remove the `length.thrift` from the `/idl` subdirectory, and clean up the rest of the generated server code by deleting all the files.
Perform a test to the deleted `length` service:
```shell
curl -X POST http://localhost:8080/gateway/length/LengthMethod \
    -H "Content-Type: application/json" \
    -d '{"msg":"hello"}' \
    -w '\n'
```
You should see:
```shell
{"err_code":10002,"err_message":"server not found"}
```
