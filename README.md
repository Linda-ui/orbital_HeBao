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

For Milestone 2, we have implemented the API gateway plus two backend Kitex services. For each Kitex service, we created three duplicate instances (hosted on different ports) for testing purposes. 

The features we implemented includes:
1. **Service registration and discovery:** Nacos is integrated in our gateway as the service registry. Kitex service instances register their addresses on the Nacos server, and Kitex clients can discover these instances through the Nacos server to establish connections with the corresponding servers.

2. **Load balacing:** Kitex's default load balancer is integrated in our gateway so that loads (request traffics) can be evenly distributed on the three server instances created for each Kitex service. This feature is to prevent overloading and thus breakdown of any single service instance.

3. **IDL-service mapping:** our gateway maintains a dynamic map between each interface definition language (IDL) file and its corresponding Kitex service. When an IDL file describing a service is added/deleted, the map will update itself to reflect changes in the managed services of our gateway.
<br></br>

## Initialisation

First [download the latest version of Nacos](https://github.com/alibaba/nacos/releases)

Then, run the gateway startup script to start up the API gateway and the nacos server locally.
```shell
# run at project root directory
sh script/gateway_startup_zsh.sh # for zshell
# OR
bash script/gateway_startup_bash.sh # for bash
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
sh script/server_startup.sh
```
<br>

## Sending Request

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

## Unit Test
Unit testing for the gateway is partially done. Run the following commands to run all unit tests and check test coverage:
```shell
go test ./...
go test ./... -cover
```