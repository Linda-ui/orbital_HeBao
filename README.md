# API Gateway based on CloudWeGo Projects
This is the orbital23 project for team HeBao.
<br></br>

## Proposed Level of Achievement
Artemis
<br></br>

## Motivation
In recent years, microservice architecture has experienced a significant increase in popularity due to its flexibility and scalability. Many companies such as ByteDance, Google, and Netflix are leveraging on microservices architecture to deliver a more seamless user experience. Many aspiring small-scale companies are also riding the tide of microservices for better managing larger-scale applications with large number of microservices. Yet, managing client-server interactions with numerous backend microservices can be complex. To address this challenge, we propose implementing an API gateway as the single point of communication between the client and multiple backend microservices. 

The gateway decouples clients from individual microservices, enabling independent service evolution while providing a unified interface for clients. We will implement logic like service registration and discovery, load balancing, rate limiting, and more for our API gateway to distribute traffic and streamline the client-server communication process. 

Eventually, we hope to build an API gateway that simplifies the management of client-server interactions and enhances the efficiency and performance of the microservices architecture.
<br></br>

## Project Scope
The scope of our project is to develop an API gateway using CloudWego's Kitex (an RPC framework) and Hertz (an HTTP framework) for managing and securing communication between clients and microservices.

To elaborate, our project aims to create an API gateway solution based on CloudWego's Kitex and Hertz framework, which will serve as a central entry point for clients to access multiple microservices. The API gateway will handle tasks such as request routing, load balancing, rate limiting, service registry and discovery, parameter binding and validation, as well as providing analytics and reporting capabilities. The goal is to enhance the overall performance, scalability, and security of the system architecture by consolidating and managing communication between clients and the underlying microservices through a unified and streamlined interface.
<br></br>

## Overview
The API gateway is mainly built with two Golang frameworks, Kitex and Hertz. The most basic function of the gateway is to translate the HTTP request from client to Thrift binary protocol, conduct a remote procedure call (RPC) to the backend RPC servers, and finally translate the response back to HTTP. 

Various extensions are integrated into the API gateway to handle extra logic besides the basic function. We use nacos (an extension for Kitex) for service registration and discovery. Thrift Interface Definition Language (IDL) is also used to define the interface contracts between the client, the gateway, and backend microservices. 
<br></br>

## System Design
![system design](./docs/System%20Design%20Diagram.jpg)
The API gateway is implemented as a Hertz server containing multiple Kitex clients, each corresponding to a microservice (implemented as a Kitex server). The Hertz
server receives incoming HTTP requests from the client and invokes the respective Kitex client. Then, the Kitex client will translate the HTTP request to Thrift binary request and initiate a generic RPC call to its backend server. The backend Kitex server handles the request and sends a response back to the Kitex client, which is then translated back to HTTP response and directed back to the client. 
<br></br>

## Components
### Parameter Binding and Validation

This component will process the HTTP request provided by the client for correct identification, extraction and validation of information from the request to be used in the application.

This is generated using the Hertz framework.

### HTTP Mapping Generic Call
This component converts the HTTP request into a generic request based on the interface mapping specifications in the Thrift IDL files. The generic call is then performed to obtain a generic response which is converted back into HTTP response.

This is implemented using Kitex’s generic call feature.

### Service Registry and Discovery
This component refers to a service registry that is responsible for registering and identifying services the API gateway interacts with. The gateway dynamically discovers the service location through the registry, establishes connection and communicates with the service needed. New services can be added, removed, or updated without requiring manual reconfiguration of the API gateway.

This is implemented using an extension for Kitex, nacos.

### Load Balancing
This component distributes incoming API requests across multiple instances of backend services to ensure even distribution of workload and optimal resource utilization. This component helps achieve scalability by adding or removing backend service instances dynamically based on the incoming traffic.

This is implemented using Kitex’s default load-balancer, which uses a weighted round-robin strategy based on weights.

### Rate Limiting
This component prevents the server from being overloaded by sudden traffic increase from a client by controlling the rate of incoming requests for each rate-limited endpoint.

This is implemented using Kitex’s default rate limiter, which specifies the maximum number of concurrent connections and the maximum number of queries per second to a specific endpoint.

### Interactions between Components
1. The Hertz Server component is responsible for receiving incoming requests from clients. It performs parameter binding and validation, extracting data from the request and mapping it to the corresponding API parameters. The Hertz Server then passes the processed request to the generic client.
2. The generic client analyzes the request metadata, such as the HTTP method, headers, URL path, and query parameters, to map the HTTP request to a generic call to the corresponding backend service.
3. The incoming request first passes through the rate limiter component in the API gateway. The rate limiter evaluates the request against predefined limits, such as the number of requests per second (QPS) or concurrent connections. If the request exceeds the configured limits, it may be delayed or rejected based on the defined behavior.
4. If the request is accepted, nacos will help us discover the available instances of the backend service required to process the request from all services registered with it.
5. The load balancer uses a predefined algorithm, such as weighted round-robin, to select an instance from the pool of available instances to evenly distribute the incoming requests across the backend service instances to achieve optimal resource utilization and performance.
6. The request is forwarded to the selected instance of the microservice for process and a response will be sent back to the generic client.
7. The generic response is translated back into HTTP response and relayed back to the client. 
<br></br>

## Milstone 1 Technical Proof
For Milestone 1, we have implemented the API gateway plus two simple backend Kitex servers for testing. We also incorporated Nacos as a service registry centre for more efficient discovery of our backend services. Load balancing and rate limiting features are not implemented yet.

### Initialisation

First download Nacos and start the Nacos server locally. The link for download is https://github.com/alibaba/nacos/releases and the latest version is used for testing. 

To check if the Nacos server is running, log in at http://127.0.0.1:8848/nacos/index.html#/login with username `nacos` and passwor `nacos`.

Then, run the API gateway in shell.
```shell
# root directory
go run ./hertz_gateway 
```

Then, run the two Kitex test servers.
```shell
# kitex_servers/server
go run ./echo
go run ./sum
```
### Testing

test if the API gateway is running
```shell
curl http://localhost:8080/

# "the api gateway is running"% 
```

test the `echo` service
```shell
curl -X POST http://localhost:8080/gateway/echo -H "Content-Type: application/json" -d '{"method":"echomethod","biz_params":"{\"msg\":\"hello\"}"}'
```
```shell
{"err_code":0,"err_message":"ok","msg":"hello"}
```

test the `sum` service
```shell
curl -X POST http://localhost:8080/gateway/sum -H "Content-Type: application/json" -d '{"method":"summethod","biz_params":"{\"firstNum\":2,\"secondNum\":4}"}'
```
```shell
{"err_code":0,"err_message":"ok","sum":6}
```

## Future Development
We propose the following improvements to be made in the future:
- Enhance the security features:
Implement features such as OAuth 2.0 to have more robust authentication and authorization mechanisms to protect APIs and sensitive data.
- Implement dynamic configuration: 
Allow for on-the-fly changes to routing rules, rate limits, security policies, and other parameters. This enables more flexible and agile management of the gateway configuration without requiring a complete restart or redeployment.
- Prepare for deployment:
Make our API gateway available for deployment by users. Several steps include: Package our API gateway software and related configuration files into a deployable artifact. Provide comprehensive documentation that guides clients through the deployment process. Create automation scripts or configuration templates that streamline the deployment process, with tools such as Docker Compose.
<br></br>

## Conclusion
Our API gateway will serve as a powerful interface that decouples the client interface from backend implementation and manage all microservices in one place. It processes all incoming requests with various functions like rate-limiting and load-balancing. Our gateway will be especially helpful for hosting large-scale APIs and it greatly improves the flexibility and performance of the entire application.
