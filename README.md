# API Gateway based on CloudWeGo Projects
This is the orbital23 project for team HeBao.
<br></br>


## Proposed Level of Achievement
Artemis
<br></br>


## Documentation

For the detailed documentation of this project including design, features, and development logs can be accessed at [this page](https://cloud-orchid-e5c.notion.site/API-Gateway-based-on-CloudWeGo-Projects-6b6f65e1a3034a8d8d1a98af719a884a?pvs=4).
<br></br>


## Overview

For Milestone 3, we have implemented the API gateway plus two backend Kitex services. For each Kitex service, we created three duplicate instances (hosted on different ports) for testing purposes. 

The features we implemented includes:
1. **Service registration and discovery:** Nacos is integrated in our gateway as the service registry. Kitex service instances register their addresses on the Nacos server, and Kitex clients can discover these instances through the Nacos server to establish connections with the corresponding servers.

2. **Load balacing:** Kitex's default load balancer is integrated in our gateway so that loads (request traffics) can be evenly distributed on the three server instances created for each Kitex service. This feature is to prevent overloading and thus breakdown of any single service instance.

3. **IDL-service mapping and dynamic update:** our gateway maintains a dynamic map between each interface definition language (IDL) file and its corresponding Kitex service. When an IDL file describing a service is added/deleted, the map will update itself to reflect changes in the managed services of our gateway, even when the gateway is still running. The client would not face disruptions of the gateway service when the backend services are being added or deleted.
<br></br>

## Initialisation

First [download the latest version of Nacos](https://github.com/alibaba/nacos/releases)

Then, run the gateway startup script to start up the API gateway and the nacos server locally.
```shell
# run at project root directory
./script/gateway_startup_zsh.sh # for zshell
# OR
./script/gateway_startup_bash.sh # for bash
```

To check if the API gateway is running.
```shell
curl http://localhost:8080/

# should return "the api gateway is running"
```

To check if the Nacos server is running, log in at http://127.0.0.1:8848/nacos/index.html#/login with username `nacos` and password `nacos`.

Finally, run the backend server startup script to start up all service instances.
```shell
# run at project root directory
./script/server_startup.sh
```
<br>

## Integration test

To use the `echo` service, run the script that sends various valid/invalid requests to the gateway as follow:
```shell
# run at project root directory
sh script/echo_request.sh
```

The gateway is customised to handle business errors like invalid user requests. Depending on the user inputs, various error messages will be returned with corresponding error codes (in the form of a response). This is the output for running `echo_request.sh`:
```shell
# Normal request 1
{"err_code":0,"err_message":"success","msg":"hello"}
# Normal request 2
{"err_code":0,"err_message":"success","msg":"goodbye"}
# Bad request 1: request with INVALID json body
{"err_code":10001,"err_message":"bad request"}
# Bad request 2: request with EMPTY json body
{"err_code":10001,"err_message":"bad request"}
# Server not found. The server echoXXX does not exist.
{"err_code":10002,"err_message":"server not found"}
# Server method not found. The method EchoMethodXXX does not exist.
{"err_code":10003,"err_message":"server method not found"}
```

Similarly, to send requests to the `sum` service, run
```shell
# run at project root directory
sh script/sum_request.sh
```

The output is
```shell
# Normal request 1: both positive numbers
{"err_code":0,"err_message":"success","sum":6}
# Normal request 2: with a negative number
{"err_code":0,"err_message":"success","sum":-6000}
```

Testing for business errors is omitted since it is similar to that of `echo` service.
<br></br>

## Integration test
Navigate to the `/integration_test` directory, and run the go tests in verbose mode to see all our integration testcases:
```shell
go test ./... -v
```

## Unit Test
Several chosen packages have unit tests done. Run the following commands to run all unit tests and check test coverage:
```shell
go test ./...
go test ./... -cover
```

## Testing the dynamic update feature during runtime of the gateway
First, let us start running the gateway and example servers `echo` and `sum`.
### Adding a service
The official documentation for creating a new Kitex service can be found [here](https://www.cloudwego.io/docs/kitex/getting-started/). You can also follow our instructions below for a simplified version of a service for testing. If you face any difficulties, please refer to the official doc for debugging.

Let us add an example `length` Kitex service that returns the length of the string input. You can also use another service with a different service name. 

Under our project root directory, create a new Thrift file that defines this service and is used for Kitex to generate skeleton code.
By convention, we name it `length.thrift`.
```thrift
namespace go length

struct LengthReq {
  1: string msg
}

struct LengthResp {
  1: i64 msg
}

service LengthSvc {
  LengthResp LengthMethod(1: LengthReq req)
}
```
You can also add in other structs and service methods as needed.
After saving this file, we use the Kitex compiler to compile the IDL file to generate whole project.
```
$ kitex -module github.com/Linda-ui/orbital_HeBao -service length length.thrift
```
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
You  will notice that there is not a subdirectory called `handler`. Please create that subdirectory yourself and put `handler.go` inside. Meanwhile, you should also change the package name of `handler.go` to be `handler`. You can do this by changing on first line of `handler.go`.

For cleaner code structure, you can place all the above generated file into a new subdirectory `length` under `kitex_services/`. Note that you may have to rewrite the import path of package `length` in all generated files with this import package. Skipping this step will not affect the correctness of code. For ease of illustration, I will write this service in the root directory.

All method process entry should be in `handler.go`. Below is an example implementation of the length service.
```
package handler

import (
	"context"
	length "github.com/Linda-ui/orbital_HeBao/kitex_gen/length"
)

// LengthSvcImpl implements the last service interface defined in the IDL.
type LengthSvcImpl struct{}

// LengthMethod implements the LengthSvcImpl interface.
func (s *LengthSvcImpl) LengthMethod(ctx context.Context, req *length.LengthReq) (resp *length.LengthResp, err error) {
	return &length.LengthResp{Msg: int64(len(req.Msg))}, nil
}
```

Next, modify the `main.go` file so that it registers the service with Nacos, has a service name called "length", and uses a customised port number 8893. You may change this port number but you should avoid ports already used by the gateway and other services, or any other programs currently running on your laptop. By default, the gateway runs on port 8080, the echo service runs on ports 8870, 8871, 8872, the sum service runs on ports 9870, 9871, 9872.
```
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
Next, run `main` to get the server running.
```
go run main.go
```
If you see something like this in your terminal:
```shell
2023/07/30 16:16:59.242864 logger.go:45: [Info] logDir:</home/hejin/Study/orbital_HeBao/log>   cacheDir:</home/hejin/Study/orbital_HeBao/cache>
2023/07/30 16:16:59.243099 logger.go:45: [Info] udp server start, port: 55045
2023/07/30 16:16:59.243642 server.go:81: [Info] KITEX: server listen at addr=127.0.0.1:8893
2023/07/30 16:17:00.244740 logger.go:45: [Info] register instance namespaceId:<>,serviceName:<DEFAULT_GROUP@@length> with instance:<{"valid":false,"marked":false,"instanceId":"","port":8893,"ip":"127.0.0.1","weight":10,"metadata":{},"clusterName":"DEFAULT","serviceName":"","enabled":true,"healthy":true,"ephemeral":true}>
2023/07/30 16:17:00.249035 logger.go:45: [Info] adding beat: <{"ip":"127.0.0.1","port":8893,"weight":10,"serviceName":"DEFAULT_GROUP@@length","cluster":"DEFAULT","metadata":{},"scheduled":false}> to beat map
2023/07/30 16:17:00.249103 logger.go:45: [Info] namespaceId:<> sending beat to server:<{"ip":"127.0.0.1","port":8893,"weight":10,"serviceName":"DEFAULT_GROUP@@length","cluster":"DEFAULT","metadata":{},"scheduled":false}>
```
Congratulations! You have successfully ran the new Kitex service!
Now, in order for our gateway to recognise our newly added service, we need to move `length.thrift` into `idl/` subdirectory. 

Test this new service through our gateway by:  
```shell
curl -X POST http://localhost:8080/gateway/length/LengthMethod \
    -H "Content-Type: application/json" \
    -d '{"msg":"hello"}' \
    -w '\n'
```
You should see:
```
{"err_code":0,"err_message":"success","msg":5}
```
Congratulations again! You have just added a new Kitex service to our gateway while it is still running, leveraging our dynamic update feature.

